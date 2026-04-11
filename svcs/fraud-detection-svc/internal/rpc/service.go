package rpc

import "fraud-detection-svc/internal/engine"

type Service struct {
	engine engine.ScoringEngine
}

func NewService(engine engine.ScoringEngine) *Service {
	return &Service{engine}
}

type ScorePaymentResponse struct {
	PrescriptiveAction engine.PrescriptiveAction
	Confidence         float64
	Issues             []string
}

func (s *Service) ScorePayment(tx engine.PaymentInfo) (*ScorePaymentResponse, error) {
	txRisk, err := s.engine.ScoreTx(tx)
	if err != nil {
		return nil, err
	}

	bhRisk, err := s.engine.ScoreAccount(tx.AccountID)
	if err != nil {
		return nil, err
	}

	if p, err := s.engine.Prescribe(txRisk.Score, bhRisk.Score); err != nil {
		return nil, err
	} else {
		return &ScorePaymentResponse{
			PrescriptiveAction: p.Action,
			Confidence:         p.Confidence,
			Issues:             append(txRisk.Indicators, bhRisk.Indicators...),
		}, nil
	}
}
