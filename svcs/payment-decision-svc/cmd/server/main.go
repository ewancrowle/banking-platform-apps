package main

import (
	"context"
	"database/sql"
	"fmt"
	"identity-svc/gen/identity/v1/identityv1connect"
	ledgerv1 "ledger-svc/gen/ledger/v1"
	"ledger-svc/gen/ledger/v1/ledgerv1connect"
	"log"
	"net/http"
	v1 "payment-decision-svc/gen/payment_decision/v1"
	"payment-decision-svc/gen/payment_decision/v1/payment_decisionv1connect"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/kelseyhightower/envconfig"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"google.golang.org/protobuf/types/known/emptypb"
)

type config struct {
	Port                int    `default:"8080"`
	IdentityServiceAddr string `required:"true" split_words:"true"`
	LedgerServiceAddr   string `required:"true" split_words:"true"`
	DBHost              string `envconfig:"db_host" required:"true"`
	DBName              string `envconfig:"db_name" required:"true"`
	DBUsername          string `envconfig:"db_username" required:"true"`
	DBPassword          string `envconfig:"db_password" required:"true"`
}

type service struct {
	payment_decisionv1connect.PaymentDecisionServiceHandler
	db                    *bun.DB
	identityServiceClient identityv1connect.IdentityServiceClient
	ledgerServiceClient   ledgerv1connect.LedgerServiceClient
}

func (s service) DecidePayment(ctx context.Context, request *v1.DecidePaymentRequest) (*v1.DecidePaymentResponse, error) {
	id, err := s.identityServiceClient.ID(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	balances, err := s.ledgerServiceClient.GetBalances(ctx, &ledgerv1.GetBalancesRequest{
		AccountId: request.AccountId,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if request.Type == "withdrawal" || request.Type == "card" || request.Type == "account_to_account" {
		if balances.AvailableBalance+request.Amount < 0 {
			return &v1.DecidePaymentResponse{
				Decision:      v1.Decision_DECISION_DECLINED,
				DecisionId:    id.Id,
				DeclineReason: v1.DeclineReason_DECLINE_REASON_INSUFFICIENT_FUNDS.Enum(),
			}, nil
		}
	}

	return &v1.DecidePaymentResponse{
		Decision:   v1.Decision_DECISION_APPROVED,
		DecisionId: id.Id,
	}, nil
}

func main() {
	var c config
	err := envconfig.Process("", &c)
	if err != nil {
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

	ledgerServiceClient := ledgerv1connect.NewLedgerServiceClient(
		http.DefaultClient,
		c.LedgerServiceAddr,
	)

	svc := service{
		db:                    db,
		identityServiceClient: identityServiceClient,
		ledgerServiceClient:   ledgerServiceClient,
	}

	path, handler := payment_decisionv1connect.NewPaymentDecisionServiceHandler(svc, connect.WithInterceptors(validate.NewInterceptor()))

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

	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
