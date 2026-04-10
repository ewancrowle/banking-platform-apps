package main

import (
	"account-svc/pkg/model/account"
	v1 "confirmation-of-payee-svc/gen/confirmation_of_payee/v1"
	"confirmation-of-payee-svc/gen/confirmation_of_payee/v1/confirmation_of_payeev1connect"
	"confirmation-of-payee-svc/pkg/config"
	"confirmation-of-payee-svc/pkg/model/confirmationofpayee"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"identity-svc/gen/identity/v1/identityv1connect"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/kelseyhightower/envconfig"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"google.golang.org/protobuf/types/known/emptypb"
)

type service struct {
	confirmation_of_payeev1connect.ConfirmationOfPayeeServiceHandler
	db                    *bun.DB
	identityServiceClient identityv1connect.IdentityServiceClient
}

func (s *service) ConfirmPayee(ctx context.Context, req *v1.ConfirmPayeeRequest) (*v1.ConfirmPayeeResponse, error) {
	a := new(account.Account)
	exists, err := s.db.NewSelect().
		Model(a).
		Where("account_num = ?", req.AccountNum).
		Where("first_name ILIKE ?", req.FirstName).
		Where("last_name ILIKE ?", req.LastName).
		Exists(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if !exists {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("account not found"))
	}

	id, err := s.identityServiceClient.ID(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	tokenString := base64.RawURLEncoding.EncodeToString(b)

	hash := sha256.Sum256([]byte(tokenString))
	hashString := hex.EncodeToString(hash[:])

	t := confirmationofpayee.ConfirmationOfPayeeToken{
		ID:        id.Id,
		AccountID: a.ID,
		Hash:      hashString,
	}

	if err = t.Insert(ctx, s.db); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &v1.ConfirmPayeeResponse{ConfirmationOfPayeeToken: tokenString}, nil
}

func (s *service) IntrospectToken(ctx context.Context, req *v1.IntrospectTokenRequest) (*v1.IntrospectTokenResponse, error) {
	hash := sha256.Sum256([]byte(req.ConfirmationOfPayeeToken))
	hashString := hex.EncodeToString(hash[:])

	t, err := confirmationofpayee.SelectTokenByHash(ctx, s.db, hashString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("token not found"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if !t.UsedAt.IsZero() {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("token already used"))
	}

	t.UsedAt = time.Now()
	_, err = s.db.NewUpdate().Model(t).Column("used_at").WherePK().Exec(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &v1.IntrospectTokenResponse{AccountId: t.AccountID}, nil
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

	svc := &service{
		db: db,
	}

	path, handler := confirmation_of_payeev1connect.NewConfirmationOfPayeeServiceHandler(svc, connect.WithInterceptors(validate.NewInterceptor()))

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
