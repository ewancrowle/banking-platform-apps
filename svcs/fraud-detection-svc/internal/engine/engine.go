package engine

type OperatingMode int

const (
	ModeBalanced OperatingMode = iota
	ModeStrict
	ModePermissive
)

type PrescriptiveAction int32

const (
	ActionPermit PrescriptiveAction = iota
	ActionProbe
	ActionProhibit
)

type Risk struct {
	Score      float64
	Indicators []string
}

type Prescription struct {
	Action     PrescriptiveAction
	Confidence float64
}

type PaymentInfo struct {
	PaymentID      int64
	AccountID      int64
	MerchantID     *int64
	OtherAccountID *int64
	Amount         int64
	CurrencyCode   string
	Type           string
}

type ScoringEngine interface {
	SetMode(mode OperatingMode)
	GetMode() OperatingMode
	Prescribe(txScore, bhScore float64) (Prescription, error)
	ScoreTx(tx PaymentInfo) (Risk, error)
	ScoreAccount(id int64) (Risk, error)
}
