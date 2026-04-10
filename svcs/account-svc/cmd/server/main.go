package main

import (
	v1 "account-svc/gen/account/v1"
	"account-svc/gen/account/v1/accountv1connect"
	"account-svc/internal/luhn"
	"account-svc/pkg/config"
	"account-svc/pkg/model/account"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"identity-svc/gen/identity/v1/identityv1connect"
	"log"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/alexedwards/argon2id"
	"github.com/kelseyhightower/envconfig"
	"github.com/moov-io/iso4217"
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

	num, err := luhn.Generate(8)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	hash, err := argon2id.CreateHash(request.Password, argon2id.DefaultParams)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	a := &account.Account{
		ID:           id.Id,
		AccountNum:   num,
		FirstName:    request.FirstName,
		MiddleNames:  request.GetMiddleNames(),
		LastName:     request.LastName,
		Email:        strings.ToLower(request.Email),
		PasswordHash: hash,
		Line1:        request.GetLine_1(),
		Line2:        request.Line_2,
		Town:         request.Town,
		Postcode:     request.Postcode,
		IsFrozen:     false,
		CurrencyCode: iso4217.GBP.Code,
	}

	if err := a.Insert(ctx, s.db); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &v1.CreateAccountResponse{Id: id.Id}, nil
}

func (s service) VerifyCredentials(ctx context.Context, request *v1.VerifyCredentialsRequest) (*v1.VerifyCredentialsResponse, error) {
	a, err := account.SelectByEmail(ctx, s.db, strings.ToLower(request.Email))
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

	return &v1.GetAccountResponse{
		Id:           a.ID,
		AccountNum:   a.AccountNum,
		FirstName:    a.FirstName,
		MiddleNames:  a.MiddleNames,
		LastName:     a.LastName,
		Email:        a.Email,
		Line_1:       a.Line1,
		Line_2:       a.Line2,
		Town:         a.Town,
		Postcode:     a.Postcode,
		IsFrozen:     a.IsFrozen,
		CurrencyCode: a.CurrencyCode,
		CreatedAt:    a.CreatedAt.String(),
	}, nil
}

func main() {
	var c config.Config
	if err := envconfig.Process("", &c); err != nil {
		log.Fatal(err.Error())
	}

	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(c.DBHost),
		pgdriver.WithDatabase(c.DBName),
		pgdriver.WithUser(c.DBUsername),
		pgdriver.WithPassword(c.DBPassword),
		pgdriver.WithInsecure(true),
	))

	db := bun.NewDB(sqlDB, pgdialect.New()).WithQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(true),
		bundebug.FromEnv(),
	))

	identityServiceClient := identityv1connect.NewIdentityServiceClient(
		http.DefaultClient,
		c.IdentityServiceAddr,
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
		Addr:      fmt.Sprintf(":%d", c.Port),
		Handler:   mux,
		Protocols: p,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
