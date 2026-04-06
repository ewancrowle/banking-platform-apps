package main

import (
	accountv1 "account-svc/gen/account/v1"
	"account-svc/gen/account/v1/accountv1connect"
	"context"
	"database/sql"
	"errors"
	"fmt"
	v1 "ledger-svc/gen/ledger/v1"
	"ledger-svc/gen/ledger/v1/ledgerv1connect"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"github.com/kelseyhightower/envconfig"
	"github.com/uptrace/go-clickhouse/ch"
	"github.com/uptrace/go-clickhouse/chdebug"
)

type config struct {
	Port               int    `default:"8080"`
	AccountServiceAddr string `required:"true" split_words:"true"`
	ClickHouseURL      string `envconfig:"clickhouse_url" required:"true"`
}

type service struct {
	ledgerv1connect.LedgerServiceHandler
	accountServiceClient accountv1connect.AccountServiceClient
	chDB                 *ch.DB
}

func (s service) GetBalances(ctx context.Context, request *v1.GetBalancesRequest) (*v1.GetBalancesResponse, error) {
	a, err := s.accountServiceClient.GetAccount(ctx, &accountv1.GetAccountRequest{
		Id: request.AccountId,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("account not found"))
	}

	var capturedAmount int64
	err = s.chDB.NewSelect().
		ColumnExpr("sum(total_amount) AS total_amount").
		Table("total_amount_captured").
		Where("account_id = ?", request.AccountId).
		Where("currency_code = ?", a.CurrencyCode).
		Scan(ctx, &capturedAmount)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var pendingAmount int64
	inner := s.chDB.NewSelect().
		ColumnExpr("sumMerge(authorised_amount) + sumMerge(incremented_amount) AS pending_amount").
		Table("pending_payments").
		Where("account_id = ?", request.AccountId).
		Where("currency_code = ?", a.CurrencyCode).
		Group("id").
		Having("maxMerge(is_captured) = 0").
		String()
	err = s.chDB.NewRaw("SELECT sum(pending_amount) AS total_amount FROM ("+inner+")").Scan(ctx, &pendingAmount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &v1.GetBalancesResponse{
		CurrentBalance:   capturedAmount,
		AvailableBalance: capturedAmount + pendingAmount,
		CurrencyCode:     a.CurrencyCode,
	}, nil
}

func (s service) GetTotalSpending(ctx context.Context, request *v1.GetTotalSpendingRequest) (*v1.GetTotalSpendingResponse, error) {
	_, err := s.accountServiceClient.GetAccount(ctx, &accountv1.GetAccountRequest{
		Id: request.AccountId,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("account not found"))
	}

	//TODO implement me
	panic("implement me")
}

func main() {
	var c config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	chDB := ch.Connect(ch.WithDSN(c.ClickHouseURL))
	chDB.AddQueryHook(chdebug.NewQueryHook(
		chdebug.WithEnabled(false),
		chdebug.FromEnv(),
	))

	accountServiceClient := accountv1connect.NewAccountServiceClient(
		http.DefaultClient,
		c.AccountServiceAddr,
	)

	svc := service{
		accountServiceClient: accountServiceClient,
		chDB:                 chDB,
	}

	path, handler := ledgerv1connect.NewLedgerServiceHandler(svc)

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
