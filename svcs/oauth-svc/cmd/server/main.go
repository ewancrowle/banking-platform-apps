package main

import (
	"context"
	"net/http"
	v1 "oauth-svc/gen/oauth/v1"
	"oauth-svc/gen/oauth/v1/oauthv1connect"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
)

type service struct {
	oauthv1connect.OAuthServiceHandler
}

func (s service) Token(ctx context.Context, request *v1.TokenRequest) (*v1.TokenResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) Refresh(ctx context.Context, request *v1.RefreshRequest) (*v1.TokenResponse, error) {
	//TODO implement me
	panic("implement me")
}

func main() {
	svc := service{}

	path, handler := oauthv1connect.NewOAuthServiceHandler(svc, connect.WithInterceptors(validate.NewInterceptor()))

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
