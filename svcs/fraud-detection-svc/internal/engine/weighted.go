package engine

import "sync"

type engineParams struct {
	TxWeight           float64 // Weighting for transactional indicators
	BehaviourWeight    float64 // Weighting for behavioural indicators
	ProbeLowerBound    float64 // Lower bound for probative prescriptive action
	ProhibitLowerBound float64 // Lower bound for prohibitive prescriptive action
}

// WeightedScoringEngine uses a weighted sum model to score and detect payment fraud.
type WeightedScoringEngine struct {
	mu          sync.RWMutex
	currentMode OperatingMode
	modes       map[OperatingMode]engineParams
}

func NewWeightedScoringEngine(initialMode OperatingMode) *WeightedScoringEngine {
	return &WeightedScoringEngine{
		currentMode: initialMode,
		modes: map[OperatingMode]engineParams{
			ModeBalanced: {
				TxWeight:           0.6,
				BehaviourWeight:    0.4,
				ProbeLowerBound:    0.4,
				ProhibitLowerBound: 0.7,
			},
			ModeStrict: {
				TxWeight:           0.8,
				BehaviourWeight:    0.2,
				ProbeLowerBound:    0.2,
				ProhibitLowerBound: 0.5,
			},
			ModePermissive: {
				TxWeight:           0.4,
				BehaviourWeight:    0.6,
				ProbeLowerBound:    0.6,
				ProhibitLowerBound: 0.8,
			},
		},
	}
}

// SetMode implements [ScoringEngine].
func (e *WeightedScoringEngine) SetMode(mode OperatingMode) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.currentMode = mode
}

// GetMode implements [ScoringEngine].
func (e *WeightedScoringEngine) GetMode() OperatingMode {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.currentMode
}

func (e *WeightedScoringEngine) GetParams() engineParams {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.modes[e.currentMode]
}

func invLerp(score, from, to float64) float64 {
	if to <= from {
		return 1.0
	}
	v := (score - from) / (to - from)
	if v > 1.0 {
		return 1.0
	}
	if v < 0.0 {
		return 0.0
	}
	return v
}

// Prescribe implements [ScoringEngine].
func (e *WeightedScoringEngine) Prescribe(txScore, bhScore float64) (Prescription, error) {
	p := e.GetParams()

	score := (txScore * p.TxWeight) + (bhScore * p.BehaviourWeight)

	if score < p.ProbeLowerBound {
		return Prescription{
			Action:     ActionPermit,
			Confidence: 1.0 - invLerp(score, 0.0, p.ProbeLowerBound),
		}, nil
	} else if score < p.ProhibitLowerBound {
		return Prescription{
			Action:     ActionProbe,
			Confidence: invLerp(score, p.ProbeLowerBound, p.ProhibitLowerBound),
		}, nil
	}
	return Prescription{
		Action:     ActionProhibit,
		Confidence: invLerp(score, p.ProhibitLowerBound, 1.0),
	}, nil
}

// ScoreTx implements [ScoringEngine].
func (e *WeightedScoringEngine) ScoreTx(tx PaymentInfo) (Risk, error) {
	indicators := []string{}
	score := 0.0

	// Mock rule
	if tx.Amount > 100_000_00 { // Over £100k is risky
		score += 0.25
		indicators = append(indicators, "high_amount")
	}

	return Risk{
		score,
		indicators,
	}, nil
}

// ScoreAccount implements [ScoringEngine].
func (e *WeightedScoringEngine) ScoreAccount(id int64) (Risk, error) {
	indicators := []string{}
	score := 0.0

	// Mock rule
	if true {
		score += 0.25
		indicators = append(indicators, "high_payment_velocity")
	}

	return Risk{
		score,
		indicators,
	}, nil
}
