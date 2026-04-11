package engine

import "sync"

type OperatingMode int

const (
	ModeTypical OperatingMode = iota
	ModeHeightened
	ModeRelaxed
)

type PrescriptiveAction int

const (
	ActionPermit PrescriptiveAction = iota
	ActionProbe
	ActionProhibit
)

type EngineParams struct {
	TxWeight           float32 // Weighting for transactional indicators
	BehaviourWeight    float32 // Weighting for behavioural indicators
	ProbeLowerBound    float32 // Lower bound for probative prescriptive action
	ProhibitLowerBound float32 // Lower bound for prohibitive prescriptive action
}

// Engine uses a weighted sum model to score and detect payment fraud.
type Engine struct {
	mu          sync.RWMutex
	currentMode OperatingMode
	modes       map[OperatingMode]EngineParams
}

func NewEngine(initialMode OperatingMode) *Engine {
	return &Engine{
		currentMode: initialMode,
		modes: map[OperatingMode]EngineParams{
			ModeTypical: {
				TxWeight:           0.6,
				BehaviourWeight:    0.4,
				ProbeLowerBound:    0.4,
				ProhibitLowerBound: 0.7,
			},
			ModeHeightened: {
				TxWeight:           0.8,
				BehaviourWeight:    0.2,
				ProbeLowerBound:    0.2,
				ProhibitLowerBound: 0.5,
			},
			ModeRelaxed: {
				TxWeight:           0.4,
				BehaviourWeight:    0.6,
				ProbeLowerBound:    0.6,
				ProhibitLowerBound: 0.8,
			},
		},
	}
}

func (e *Engine) GetParams() EngineParams {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.modes[e.currentMode]
}

func (e *Engine) SetMode(mode OperatingMode) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.currentMode = mode
}

func invLerp(score, from, to float32) float32 {
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

func (e *Engine) Score(txRisk, bhRisk float32) (PrescriptiveAction, float32) {
	p := e.GetParams()

	score := (txRisk * p.TxWeight) + (bhRisk * p.BehaviourWeight)

	if score < p.ProbeLowerBound {
		return ActionPermit, invLerp(score, 0.0, p.ProbeLowerBound)
	} else if score < p.ProhibitLowerBound {
		return ActionProbe, invLerp(score, p.ProbeLowerBound, p.ProhibitLowerBound)
	}
	return ActionProhibit, invLerp(score, p.ProhibitLowerBound, 1.0)
}
