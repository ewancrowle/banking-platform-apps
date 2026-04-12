package rpc

import (
	"confirmation-of-payee-svc/pkg/model/token"
	"context"
	"time"
)

type AccountInfo struct {
	FirstName string
	LastName  string
	AccountNo string
}

type Repo interface {
	CreateCOPToken(context.Context, *token.COPToken) (*token.COPToken, error)
	ReadCOPTokenByHash(context.Context, string) (*token.COPToken, error)
	UpdateCOPTokenUsedAt(context.Context, int64, time.Time) error
	ReadAccountIDByInfo(context.Context, AccountInfo) (int64, error)
}
