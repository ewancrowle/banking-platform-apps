package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"identity-svc/gen/identity/v1/identityv1connect"
	"log"
	v1 "merchant-svc/gen/merchant/v1"
	"merchant-svc/gen/merchant/v1/merchantv1connect"
	"merchant-svc/pkg/model/merchant"
	"net/http"

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
	DBHost              string `envconfig:"db_host" required:"true"`
	DBName              string `envconfig:"db_name" required:"true"`
	DBUsername          string `envconfig:"db_username" required:"true"`
	DBPassword          string `envconfig:"db_password" required:"true"`
}

type service struct {
	merchantv1connect.MerchantServiceHandler
	db                    *bun.DB
	identityServiceClient identityv1connect.IdentityServiceClient
}

func (s service) CreateMerchant(ctx context.Context, request *v1.CreateMerchantRequest) (*v1.CreateMerchantResponse, error) {
	id, err := s.identityServiceClient.ID(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	m := &merchant.Merchant{
		ID:              id.Id,
		Descriptor:      request.Descriptor_,
		ShortDescriptor: request.ShortDescriptor,
		MCC:             request.Mcc,
		Line1:           request.Line_1,
		Line2:           request.GetLine_2(),
		Town:            request.Town,
		Postcode:        request.Postcode,
		CountryCode:     request.CountryCode,
	}

	if err := m.Insert(ctx, s.db); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &v1.CreateMerchantResponse{Id: id.Id}, nil
}

func (s service) GetMerchant(ctx context.Context, request *v1.GetMerchantRequest) (*v1.Merchant, error) {
	m, err := merchant.Select(ctx, s.db, request.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("merchant not found"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &v1.Merchant{
		Id:              m.ID,
		Descriptor_:     m.Descriptor,
		ShortDescriptor: m.ShortDescriptor,
		Mcc:             m.MCC,
		Line_1:          m.Line1,
		Line_2:          m.Line2,
		Town:            m.Town,
		Postcode:        m.Postcode,
		CountryCode:     m.CountryCode,
		CreatedAt:       m.CreatedAt.String(),
		UpdatedAt:       m.UpdatedAt.String(),
	}, nil
}

func (s service) GetAllMerchants(ctx context.Context, req *emptypb.Empty) (*v1.GetAllMerchantsResponse, error) {
	if m, err := merchant.SelectAll(ctx, s.db); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	} else {
		merchants := make([]*v1.Merchant, len(m))
		for i, m := range m {
			merchants[i] = &v1.Merchant{
				Id:              m.ID,
				Descriptor_:     m.Descriptor,
				ShortDescriptor: m.ShortDescriptor,
				Mcc:             m.MCC,
				Line_1:          m.Line1,
				Line_2:          m.Line2,
				Town:            m.Town,
				Postcode:        m.Postcode,
				CountryCode:     m.CountryCode,
				CreatedAt:       m.CreatedAt.String(),
				UpdatedAt:       m.UpdatedAt.String(),
			}
		}

		return &v1.GetAllMerchantsResponse{Merchants: merchants}, nil
	}
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

	svc := service{
		db:                    db,
		identityServiceClient: identityServiceClient,
	}

	path, handler := merchantv1connect.NewMerchantServiceHandler(svc, connect.WithInterceptors(validate.NewInterceptor()))

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
