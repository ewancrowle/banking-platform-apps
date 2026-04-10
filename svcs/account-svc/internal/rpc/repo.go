package rpc

import (
	"account-svc/pkg/model/account"
	"context"
)

type Repo interface {
	CreateAccount(context.Context, *account.Account) (*account.Account, error)

	ReadAccountByID(context.Context, int64) (*account.Account, error)
	ReadAccountByEmail(context.Context, string) (*account.Account, error)

	ExistsAccountByEmail(context.Context, string) (bool, error)
}
