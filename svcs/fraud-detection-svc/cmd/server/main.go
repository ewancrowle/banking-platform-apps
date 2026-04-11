package main

import (
	"fmt"
	v1c "fraud-detection-svc/gen/fraud_detection/v1/fraud_detectionv1connect"
	"fraud-detection-svc/internal/engine"
	"fraud-detection-svc/internal/rpc"
	"fraud-detection-svc/pkg/config"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

func main() {
	var c config.Config
	if err := envconfig.Process("", &c); err != nil {
		log.Fatal(err.Error())
	}

	e := engine.NewWeightedScoringEngine(engine.ModeBalanced)
	svc := rpc.NewService(e)

	path, handler := v1c.NewFraudDetectionServiceHandler(rpc.NewHandler(svc))

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
