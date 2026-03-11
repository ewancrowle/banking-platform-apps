package main

import (
	"context"
	"fmt"
	v1 "id-service/gen/id_gen/v1"
	"id-service/gen/id_gen/v1/id_genv1connect"
	"net/http"
	"os"
	"strconv"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/sony/sonyflake/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IDGenServer struct {
	id_genv1connect.UnimplementedIDGenServiceHandler
	sf *sonyflake.Sonyflake
}

func (s IDGenServer) GenID(_ context.Context, _ *emptypb.Empty) (*v1.GenIDResponse, error) {
	id, err := s.sf.NextID()
	switch {
	case err != nil:
		return nil, err
	default:
		return &v1.GenIDResponse{Id: id}, nil
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

	idGenServer := IDGenServer{sf: sf}

	path, handler := id_genv1connect.NewIDGenServiceHandler(idGenServer, connect.WithInterceptors(validate.NewInterceptor()))

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
