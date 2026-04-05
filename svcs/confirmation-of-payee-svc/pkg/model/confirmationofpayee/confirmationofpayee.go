package confirmationofpayee

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type ConfirmationOfPayeeToken struct {
	bun.BaseModel `bun:"table:confirmation_of_payee_token"`

	ID        int64  `bun:",pk"`
	AccountID int64  `bun:",notnull"`
	Hash      string `bun:",unique,notnull"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UsedAt    time.Time `bun:",notnull"`
}

func (t *ConfirmationOfPayeeToken) Insert(ctx context.Context, db *bun.DB) error {
	_, err := db.NewInsert().Model(t).Exec(ctx)
	return err
}

func SelectTokenByHash(ctx context.Context, db *bun.DB, hash string) (*ConfirmationOfPayeeToken, error) {
	t := new(ConfirmationOfPayeeToken)
	err := db.NewSelect().Model(t).Where("hash = ?", hash).Scan(ctx)
	return t, err
}

func (t *ConfirmationOfPayeeToken) SetUsedAt(ctx context.Context, db *bun.DB, usedAt time.Time) error {
	_, err := db.NewUpdate().Model(t).Set("expires_at = ?", usedAt).WherePK().Exec(ctx)
	return err
}
