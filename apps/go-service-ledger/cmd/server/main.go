package main

import (
	accountv1 "account-svc/gen/account/v1"
	"account-svc/gen/account/v1/accountv1connect"
	"context"
	"errors"
	"fmt"
	v1 "ledger-svc/gen/ledger/v1"
	"ledger-svc/gen/ledger/v1/ledgerv1connect"
	"ledger-svc/pkg/config"
	"log"
	"net/http"
	"payment-svc/pkg/model/payment"

	"connectrpc.com/connect"
	"github.com/kelseyhightower/envconfig"
	"github.com/uptrace/go-clickhouse/ch"
	"github.com/uptrace/go-clickhouse/chdebug"
)

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

	var currentBalance int64
	if err = s.chDB.NewSelect().
		ColumnExpr("sum(amount)").
		Table("captured_payments_summed").
		Where("account_id = ?", request.AccountId).
		Where("currency_code = ?", a.CurrencyCode).
		Scan(ctx, &currentBalance); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var availableBalance int64
	if err = s.chDB.NewSelect().
		ColumnExpr("sum(amount)").
		TableExpr("current_payments FINAL").
		Where("account_id = ?", request.AccountId).
		Where("currency_code = ?", a.CurrencyCode).
		Where("status NOT IN (?)", ch.In([]payment.Status{payment.StatusReceived, payment.StatusDeclined, payment.StatusExpired, payment.StatusVoided})).
		Scan(ctx, &availableBalance); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &v1.GetBalancesResponse{
		CurrentBalance:   currentBalance,
		AvailableBalance: availableBalance,
		CurrencyCode:     a.CurrencyCode,
	}, nil
}

func (s service) GetTotalSpending(ctx context.Context, request *v1.GetTotalSpendingRequest) (*v1.GetTotalSpendingResponse, error) {
	a, err := s.accountServiceClient.GetAccount(ctx, &accountv1.GetAccountRequest{
		Id: request.AccountId,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("account not found"))
	}

	type spending struct {
		Today     int64 `ch:"today"`
		ThisWeek  int64 `ch:"this_week"`
		ThisMonth int64 `ch:"this_month"`
	}

	var sp spending
	if err = s.chDB.NewSelect().
		ColumnExpr("sumIf(amount, event_date = today()) AS today").
		ColumnExpr("sumIf(amount, event_date >= toStartOfWeek(today())) AS this_week").
		ColumnExpr("sumIf(amount, event_date >= toStartOfMonth(today())) AS this_month").
		TableExpr("daily_outgoing_payments FINAL").
		Where("account_id = ?", request.AccountId).
		Where("currency_code = ?", a.CurrencyCode).
		Where("status IN (?)", ch.In([]string{"authorised", "captured"})).
		Scan(ctx, &sp); err != nil {
		return nil, err
	}

	return &v1.GetTotalSpendingResponse{
		TotalSpentToday:     -sp.Today,
		TotalSpentThisWeek:  -sp.ThisWeek,
		TotalSpentThisMonth: -sp.ThisMonth,
		CurrencyCode:        a.CurrencyCode,
	}, nil
}

func main() {
	var c config.Config
	if err := envconfig.Process("", &c); err != nil {
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

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
