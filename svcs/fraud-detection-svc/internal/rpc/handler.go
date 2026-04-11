package rpc

import (
	"context"
	v1 "fraud-detection-svc/gen/fraud_detection/v1"
	"fraud-detection-svc/internal/engine"
)

type Handler struct {
	svc *Service
}

// ScorePayment implements [fraud_detectionv1connect.FraudDetectionServiceHandler].
func (h *Handler) ScorePayment(ctx context.Context, req *v1.ScorePaymentRequest) (*v1.ScorePaymentResponse, error) {
	if res, err := h.svc.ScorePayment(engine.PaymentInfo{
		PaymentID:      req.PaymentId,
		AccountID:      req.AccountId,
		MerchantID:     req.MerchantId,
		OtherAccountID: req.OtherAccountId,
		Amount:         req.Amount,
		CurrencyCode:   req.CurrencyCode,
		Type:           req.Type,
	}); err != nil {
		return nil, err
	} else {
		return &v1.ScorePaymentResponse{
			PrescriptiveAction: v1.PrescriptiveAction(res.PrescriptiveAction),
			Confidence:         res.Confidence,
			Issues:             res.Issues,
		}, nil
	}
}

// SetOperatingMode implements [fraud_detectionv1connect.FraudDetectionServiceHandler].
func (h *Handler) SetOperatingMode(ctx context.Context, req *v1.SetOperatingModeRequest) (*v1.SetOperatingModeResponse, error) {
	h.svc.engine.SetMode(engine.OperatingMode(req.Mode))
	return &v1.SetOperatingModeResponse{
		Mode: v1.OperatingMode(h.svc.engine.GetMode()),
	}, nil
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc}
}
