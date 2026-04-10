package main

import (
	v1c "confirmation-of-payee-svc/gen/confirmation_of_payee/v1/confirmation_of_payeev1connect"
	"confirmation-of-payee-svc/internal/rpc"
	"confirmation-of-payee-svc/pkg/config"
	"database/sql"
	"fmt"
	identityv1c "identity-svc/gen/identity/v1/identityv1connect"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/kelseyhightower/envconfig"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

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

	identitySvc := identityv1c.NewIdentityServiceClient(
		http.DefaultClient,
		c.IdentityServiceAddr,
	)

	repo := rpc.NewBunRepo(db)
	svc := rpc.NewService(repo, identitySvc)

	path, handler := v1c.NewConfirmationOfPayeeServiceHandler(rpc.NewHandler(svc), connect.WithInterceptors(validate.NewInterceptor()))

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
