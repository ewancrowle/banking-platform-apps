package account

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Account struct {
	bun.BaseModel `bun:"table:accounts"`

	ID           int64  `bun:",pk"`
	FirstName    string `bun:",notnull"`
	MiddleNames  string `bun:",notnull"`
	LastName     string `bun:",notnull"`
	Email        string `bun:",unique,notnull"`
	PasswordHash string `bun:",notnull"`
	IsFrozen     bool   `bun:",notnull"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	ClosedAt  time.Time `bun:",soft_delete,nullzero"`
}

func (a *Account) Insert(ctx context.Context, db *bun.DB) error {
	_, err := db.NewInsert().Model(a).Exec(ctx)
	return err
}

func Select(ctx context.Context, db *bun.DB, id int64) (*Account, error) {
	a := new(Account)
	err := db.NewSelect().Model(a).Where("id = ?", id).Scan(ctx)
	return a, err
}

func SelectByEmail(ctx context.Context, db *bun.DB, email string) (*Account, error) {
	a := new(Account)
	err := db.NewSelect().Model(a).Where("email = ?", email).Scan(ctx)
	return a, err
}

func (a *Account) Update(ctx context.Context, db *bun.DB) error {
	_, err := db.NewUpdate().Model(a).WherePK().Exec(ctx)
	return err
}

func (a *Account) Delete(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDelete().Model(a).WherePK().Exec(ctx)
	return err
}

func Exists(ctx context.Context, db *bun.DB, email string) (bool, error) {
	return db.NewSelect().Model((*Account)(nil)).Where("email = ?", email).Exists(ctx)
}
