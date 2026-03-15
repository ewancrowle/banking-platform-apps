package token

import (
	"account-svc/pkg/model/account"
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Token struct {
	ID        int64 `bun:",pk"`
	IsRevoked bool  `bun:"is_revoked,notnull"`

	AccountID *int64           `bun:"account_id,notnull"`
	Account   *account.Account `bun:"rel:belongs-to,join:account_id=id"`

	DeviceID *int64 `bun:"device_id,notnull"`
	//Device   *device.Device `bun:"rel:belongs-to,join:device_id=id"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	ExpiresAt time.Time `bun:"expires_at,nullzero,notnull"`
}

type AccessToken struct {
	Token
	bun.BaseModel `bun:"table:access_tokens"`
}

type RefreshToken struct {
	Token
	bun.BaseModel `bun:"table:refresh_tokens"`

	Hash string `bun:",unique,notnull"`
}

func (t *AccessToken) Insert(ctx context.Context, db *bun.DB) error {
	_, err := db.NewInsert().Model(t).Exec(ctx)
	return err
}

func SelectAccessToken(ctx context.Context, db *bun.DB, id int64) (*AccessToken, error) {
	t := new(AccessToken)
	err := db.NewSelect().Model(t).Relation("Account").Where("id = ?", id).Scan(ctx)
	return t, err
}

func (t *AccessToken) Delete(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDelete().Model(t).WherePK().Exec(ctx)
	return err
}

func (t *RefreshToken) Insert(ctx context.Context, db *bun.DB) error {
	_, err := db.NewInsert().Model(t).Exec(ctx)
	return err
}

func SelectRefreshTokenByHash(ctx context.Context, db *bun.DB, hash string) (*RefreshToken, error) {
	t := new(RefreshToken)
	err := db.NewSelect().Model(t).Where("hash = ?", hash).Scan(ctx)
	return t, err
}

func (t *RefreshToken) Delete(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDelete().Model(t).WherePK().Exec(ctx)
	return err
}
