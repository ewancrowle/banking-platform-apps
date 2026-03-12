package main

import (
	v1 "account-service/gen/account/v1"
	"account-service/gen/account/v1/accountv1connect"
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
)

type service struct {
	accountv1connect.UnimplementedAccountServiceHandler
}

func (s service) CreateAccount(context.Context, *v1.CreateAccountRequest) (*v1.CreateAccountResponse, error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("account.v1.AccountService.CreateAccount is not implemented"))
}

func main() {
	svc := service{}

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
