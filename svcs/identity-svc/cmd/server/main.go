package main

import (
	"context"
	"fmt"
	v1 "identity-svc/gen/identity/v1"
	"identity-svc/gen/identity/v1/identityv1connect"
	"identity-svc/pkg/config"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/sony/sonyflake/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

type service struct {
	identityv1connect.IdentityServiceHandler
	sf *sonyflake.Sonyflake
}

func (s service) ID(_ context.Context, _ *emptypb.Empty) (*v1.IDResponse, error) {
	if id, err := s.sf.NextID(); err != nil {
		return nil, err
	} else {
		return &v1.IDResponse{Id: id}, nil
	}
}

func main() {
	var c config.Config
	if err := envconfig.Process("", &c); err != nil {
		log.Fatal(err.Error())
	}

	sf, err := sonyflake.New(sonyflake.Settings{MachineID: func() (int, error) {
		return c.MachineID, nil
	}})
	if err != nil {
		panic(err)
	}

	svc := service{sf: sf}

	path, handler := identityv1connect.NewIdentityServiceHandler(svc)

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
