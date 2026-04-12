package rpc

import (
	"account-svc/pkg/model/account"
	"confirmation-of-payee-svc/pkg/model/token"
	"context"
	"time"

	"github.com/uptrace/bun"
)

type BunRepo struct {
	db *bun.DB
}

// CreateCOPToken creates a new Confirmation of Payee token in the database.
func (r *BunRepo) CreateCOPToken(ctx context.Context, t *token.COPToken) (*token.COPToken, error) {
	_, err := r.db.NewInsert().Model(t).Exec(ctx)
	return t, err
}

// ReadCOPTokenByHash finds a Confirmation of Payee token by its hash.
func (r *BunRepo) ReadCOPTokenByHash(ctx context.Context, hash string) (*token.COPToken, error) {
	t := new(token.COPToken)
	err := r.db.NewSelect().Model(t).Where("hash = ?", hash).Scan(ctx)
	return t, err
}

// UpdateCOPTokenUsedAt updates the used_at field of a Confirmation of Payee token.
func (r *BunRepo) UpdateCOPTokenUsedAt(ctx context.Context, tokenID int64, usedAt time.Time) error {
	_, err := r.db.NewUpdate().
		Model((*token.COPToken)(nil)).Set("used_at = ?", usedAt).Where("id = ?", tokenID).Exec(ctx)
	return err
}

// ReadAccountIDByInfo finds the account ID by the given account info.
func (r *BunRepo) ReadAccountIDByInfo(ctx context.Context, info AccountInfo) (int64, error) {
	var id int64
	err := r.db.NewSelect().
		Model((*account.Account)(nil)).
		Column("id").
		Where("first_name ILIKE ?", info.FirstName).
		Where("last_name ILIKE ?", info.LastName).
		Where("account_no = ?", info.AccountNo).
		Scan(ctx, &id)
	return id, err
}

func NewBunRepo(db *bun.DB) *BunRepo {
	return &BunRepo{db}
}
