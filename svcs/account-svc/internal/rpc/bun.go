package rpc

import (
	"account-svc/pkg/model/account"
	"context"

	"github.com/uptrace/bun"
)

type BunRepo struct {
	// Repo
	db *bun.DB
}

// NewBunRepo creates a new BunRepo with the given database connection.
func NewBunRepo(db *bun.DB) *BunRepo {
	return &BunRepo{db: db}
}

// CreateAccount creates a new account in the database.
func (r *BunRepo) CreateAccount(ctx context.Context, account *account.Account) (*account.Account, error) {
	_, err := r.db.NewInsert().Model(account).Exec(ctx)
	return account, err
}

// ReadAccountByID finds an account by its ID.
func (r *BunRepo) ReadAccountByID(ctx context.Context, id int64) (*account.Account, error) {
	a := new(account.Account) // Pointer to the model to scan into
	err := r.db.NewSelect().Model(a).Where("id = ?", id).Scan(ctx)
	return a, err
}

// ReadAccountByEmail finds an account by its email address.
func (r *BunRepo) ReadAccountByEmail(ctx context.Context, email string) (*account.Account, error) {
	a := new(account.Account) // Pointer to the model to scan into
	err := r.db.NewSelect().Model(a).Where("email = ?", email).Scan(ctx)
	return a, err
}

// ExistsAccountByEmail checks if an account with the given email exists.
func (r *BunRepo) ExistsAccountByEmail(ctx context.Context, email string) (bool, error) {
	return r.db.NewSelect().Model((*account.Account)(nil)).Where("email = ?", email).Exists(ctx)
}
