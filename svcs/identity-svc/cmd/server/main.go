package main

import (
	"context"
	"fmt"
	v1 "identity-svc/gen/identity/v1"
	"identity-svc/gen/identity/v1/identityv1connect"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/kelseyhightower/envconfig"
	"github.com/sony/sonyflake/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

type config struct {
	Port      int `default:"8080"`
	MachineID int `envconfig:"machine_id" required:"true"`
}

type service struct {
	identityv1connect.IdentityServiceHandler
	sf *sonyflake.Sonyflake
}

func (s service) ID(_ context.Context, _ *emptypb.Empty) (*v1.IDResponse, error) {
	id, err := s.sf.NextID()
	switch {
	case err != nil:
		return nil, err
	default:
		return &v1.IDResponse{Id: id}, nil
	}
}

func main() {
	var c config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	sf, err := sonyflake.New(sonyflake.Settings{MachineID: func() (int, error) {
		return c.MachineID, nil
	}})
	if err != nil {
		panic(err)
	}

	svc := service{sf: sf}

	path, handler := identityv1connect.NewIdentityServiceHandler(svc, connect.WithInterceptors(validate.NewInterceptor()))

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
