package payment

import (
	"account-svc/pkg/model/account"
	"context"
	"errors"
	"merchant-svc/pkg/model/merchant"
	"time"

	"github.com/uptrace/bun"
)

type (
	Status string
	Type   string
)

const (
	StatusReceived       Status = "received"
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
	ID            int64 `bun:",pk"`
	AccountID     int64 `bun:",notnull"`

	PaymentID *int64
	Payments  []Payment `bun:"rel:has-many,join:id=payment_id" json:"-"`

	MerchantID *int64
	Merchant   *merchant.Merchant `bun:"rel:has-one,join:merchant_id=id" json:"-"`

	OtherAccountID *int64
	OtherAccount   account.Account `bun:"rel:has-one,join:other_account_id=id" json:"-"`

	Amount       int64  `bun:",notnull"`
	CurrencyCode string `bun:",notnull"`

	Type        Type   `bun:",notnull"`
	Status      Status `bun:",notnull"`
	Description string `bun:",notnull"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func (p *Payment) Insert(ctx context.Context, db *bun.DB) error {
	_, err := db.NewInsert().Model(p).Exec(ctx)
	return err
}

func Select(ctx context.Context, db *bun.DB, id int64) (*Payment, error) {
	p := new(Payment)
	err := db.NewSelect().Model(p).Where("id = ?", id).Scan(ctx)
	return p, err
}
