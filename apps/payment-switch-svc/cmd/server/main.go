package main

import (
	"fmt"
	"log"
	"net/http"
	paymentv1c "payment-svc/gen/payment/v1/paymentv1connect"
	v1c "payment-switch-svc/gen/payment_switch/v1/payment_switchv1connect"
	"payment-switch-svc/internal/rpc"
	"payment-switch-svc/pkg/config"

	"github.com/kelseyhightower/envconfig"
)

func main() {
	var c config.Config
	if err := envconfig.Process("", &c); err != nil {
		log.Fatal(err.Error())
	}

	paymentSvcClient := paymentv1c.NewPaymentServiceClient(
		http.DefaultClient,
		c.PaymentServiceAddr,
	)

	path, handler := v1c.NewPaymentSwitchServiceHandler(rpc.NewHandler(paymentSvcClient))

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
