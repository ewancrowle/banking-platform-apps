package rpc

import (
	"account-svc/pkg/model/account"
	"context"
)

type Repo interface {
	CreateAccount(ctx context.Context, account *account.Account) (*account.Account, error)

	ReadAccountByID(ctx context.Context, id int64) (*account.Account, error)
	ReadAccountByEmail(ctx context.Context, email string) (*account.Account, error)

	ExistsAccountByEmail(ctx context.Context, email string) (bool, error)
}
