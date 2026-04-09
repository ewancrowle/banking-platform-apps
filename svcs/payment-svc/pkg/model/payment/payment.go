package payment

import (
	"account-svc/pkg/model/account"
	"context"
	"errors"
	"merchant-svc/pkg/model/merchant"
	"payment-decision-svc/pkg/model/paymentdecision"
	"time"

	"github.com/uptrace/bun"
)

type (
	Status string
	Type   string
)

const (
	StatusReceived       Status = "received"
	StatusDeclined       Status = "declined"
	StatusAuthorised     Status = "authorised"
	StatusIncremented    Status = "incremented"
	StatusCaptured       Status = "captured"
	StatusExpired        Status = "expired"
	StatusVoided         Status = "voided"
	TypeDeposit          Type   = "deposit"
	TypeWithdrawal       Type   = "withdrawal"
	TypeCard             Type   = "card"
	TypeAccountToAccount Type   = "account_to_account"
	TypeFee              Type   = "fee"
	TypeInterest         Type   = "interest"
	TypeRefund           Type   = "refund"
)

func GetType(s string) (Type, error) {
	t := Type(s)
	switch t {
	case TypeDeposit,
		TypeWithdrawal,
		TypeCard,
		TypeAccountToAccount,
		TypeFee,
		TypeInterest,
		TypeRefund:
		return t, nil
	}
	return "", errors.New("invalid type: " + s)
}

func (t Type) GetCorrectDirection(v int64) int64 {
	switch t {
	case TypeDeposit, TypeInterest, TypeRefund:
		if v < 0 {
			return -v
		}
		return v
	case TypeWithdrawal, TypeCard, TypeAccountToAccount, TypeFee:
		if v > 0 {
			return -v
		}
		return v
	}
	return v
}

type Payment struct {
	bun.BaseModel `bun:"table:payments"`
	ID            int64 `bun:",pk" json:"id"`
	AccountID     int64 `bun:",notnull" json:"account_id"`

	MerchantID *int64             `json:"merchant_id"`
	Merchant   *merchant.Merchant `bun:"rel:has-one,join:merchant_id=id" json:"-"`

	OtherAccountID *int64          `json:"other_account_id"`
	OtherAccount   account.Account `bun:"rel:has-one,join:other_account_id=id" json:"-"`

	Amount       int64  `bun:",notnull" json:"amount"`
	CurrencyCode string `bun:",notnull" json:"currency_code"`

	Type          Type                           `bun:",notnull" json:"type"`
	Status        Status                         `bun:",notnull" json:"status"`
	Description   string                         `bun:",notnull" json:"description"`
	DeclineReason *paymentdecision.DeclineReason `bun:",nullzero" json:"decline_reason"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
}

func (p *Payment) Insert(ctx context.Context, db *bun.DB) error {
	_, err := db.NewInsert().Model(p).Exec(ctx)
	return err
}

func Select(ctx context.Context, db *bun.DB, id int64) (*Payment, error) {
	p := new(Payment)
	err := db.NewSelect().Model(p).Where("id = ?", id).Relation("Merchant").Scan(ctx)
	return p, err
}

func SelectByAccountID(ctx context.Context, db *bun.DB, accountID int64) ([]Payment, error) {
	var payments []Payment
	err := db.NewSelect().Model(&payments).Where("account_id = ?", accountID).Relation("Merchant").Scan(ctx)
	return payments, err
}

func SelectDisplayableByAccountID(ctx context.Context, db *bun.DB, accountID int64) ([]Payment, error) {
	var payments []Payment
	err := db.NewSelect().
		Model(&payments).
		Where("account_id = ?", accountID).
		Where("status NOT IN (?)", []Status{StatusReceived, StatusExpired, StatusVoided}).
		Order("created_at DESC").
		Relation("Merchant").
		Scan(ctx)
	return payments, err
}

func (p *Payment) SetStatus(ctx context.Context, db *bun.DB, status Status) error {
	_, err := db.NewUpdate().Model(p).Set("status = ?", status).WherePK().Exec(ctx)
	return err
}

func (p *Payment) SetDeclineReason(ctx context.Context, db *bun.DB, reason paymentdecision.DeclineReason) error {
	_, err := db.NewUpdate().Model(p).Set("decline_reason = ?", reason).WherePK().Exec(ctx)
	return err
}

func (p *Payment) SetAmount(ctx context.Context, db *bun.DB, amount int64) error {
	_, err := db.NewUpdate().Model(p).Set("amount = ?", amount).WherePK().Exec(ctx)
	return err
}
