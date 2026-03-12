package main

import (
	"context"
	"fmt"
	v1 "identity-svc/gen/identity/v1"
	"identity-svc/gen/identity/v1/identityv1connect"
	"net/http"
	"os"
	"strconv"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/sony/sonyflake/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

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
	sf, err := sonyflake.New(sonyflake.Settings{MachineID: func() (int, error) {
		str := os.Getenv("MACHINE_ID")
		if str == "" {
			return 0, fmt.Errorf("MACHINE_ID env var is not set")
		}

		id, err := strconv.Atoi(str)
		if err != nil {
			return 0, fmt.Errorf("invalid MACHINE_ID: %v", err)
		}

		return id, nil
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
		Addr:      "localhost:8080",
		Handler:   mux,
		Protocols: p,
	}
	s.ListenAndServe()
}
