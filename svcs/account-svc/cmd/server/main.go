package main

import (
	v1 "account-svc/gen/account/v1"
	"account-svc/gen/account/v1/accountv1connect"
	"account-svc/pkg/model/account"
	"context"
	"database/sql"
	"errors"
	"identity-svc/gen/identity/v1/identityv1connect"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/alexedwards/argon2id"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"google.golang.org/protobuf/types/known/emptypb"
)

type service struct {
	accountv1connect.AccountServiceHandler
	db                    *bun.DB
	identityServiceClient identityv1connect.IdentityServiceClient
}

func (s service) CreateAccount(ctx context.Context, request *v1.CreateAccountRequest) (*v1.CreateAccountResponse, error) {
	exists, err := account.Exists(ctx, s.db, request.Email)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	} else if exists {
		return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("account already exists"))
	}

	id, err := s.identityServiceClient.ID(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	hash, err := argon2id.CreateHash(request.Password, argon2id.DefaultParams)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	a := &account.Account{
		ID:           id.Id,
		FirstName:    request.FirstName,
		MiddleNames:  request.MiddleNames,
		LastName:     request.LastName,
		Email:        request.Email,
		PasswordHash: hash,
	}

	if err := a.Insert(ctx, s.db); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &v1.CreateAccountResponse{Id: id.Id}, nil
}

func (s service) VerifyCredentials(ctx context.Context, request *v1.VerifyCredentialsRequest) (*v1.VerifyCredentialsResponse, error) {
	a, err := account.SelectByEmail(ctx, s.db, request.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &v1.VerifyCredentialsResponse{}, nil
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	match, err := argon2id.ComparePasswordAndHash(request.Password, a.PasswordHash)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if !match {
		return &v1.VerifyCredentialsResponse{}, nil
	}

	return &v1.VerifyCredentialsResponse{
		Id: &a.ID,
	}, nil
}

func (s service) GetAccount(ctx context.Context, request *v1.GetAccountRequest) (*v1.GetAccountResponse, error) {
	a, err := account.Select(ctx, s.db, request.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("account not found"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &v1.GetAccountResponse{Account: &v1.Account{
		FirstName:      a.FirstName,
		MiddleNames:    a.MiddleNames,
		LastName:       a.LastName,
		Email:          a.Email,
		KnownAddresses: []*v1.KnownAddress{},
	}}, nil
}

func main() {
	dsn := "postgres://postgres:@localhost:5432/test?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	db.WithQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv(),
	))

	identityServiceClient := identityv1connect.NewIdentityServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)

	svc := service{
		db:                    db,
		identityServiceClient: identityServiceClient,
	}

	path, handler := accountv1connect.NewAccountServiceHandler(svc, connect.WithInterceptors(validate.NewInterceptor()))

	mux := http.NewServeMux()
	mux.Handle(path, handler)

	p := new(http.Protocols)
	p.SetHTTP1(true)
	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)
	s := http.Server{
		Addr:      "localhost:8080",
		Handler:   mux,
		Protocols: p,
	}
	s.ListenAndServe()
}
