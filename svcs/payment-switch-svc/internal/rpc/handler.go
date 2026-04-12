package rpc

import (
	"context"
	paymentv1c "payment-svc/gen/payment/v1/paymentv1connect"
	v1 "payment-switch-svc/gen/payment_switch/v1"
	"payment-switch-svc/internal/processor"
)

type Handler struct {
	paymentSvc paymentv1c.PaymentServiceClient
}

// SwitchPayment implements [payment_switchv1connect.PaymentSwitchServiceHandler].
func (h *Handler) SwitchPayment(ctx context.Context, reqIn *v1.PaymentSwitchRequest) (*v1.PaymentSwitchResponse, error) {
	reqOut, err := processor.BuildRequest(reqIn.Payload)
	if err != nil {
		return nil, err
	}

	resAuth, err := h.paymentSvc.AuthorisePayment(ctx, reqOut)
	if err != nil {
		return nil, err
	}

	p, err := processor.BuildResponse(reqIn.Payload, resAuth)
	if err != nil {
		return nil, err
	}
	return &v1.PaymentSwitchResponse{Payload: p}, nil
}

func NewHandler(paymentSvc paymentv1c.PaymentServiceClient) *Handler {
	return &Handler{paymentSvc}
}
