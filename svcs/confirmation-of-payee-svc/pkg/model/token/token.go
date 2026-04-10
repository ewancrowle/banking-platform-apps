package token

import (
	"time"

	"github.com/uptrace/bun"
)

type COPToken struct {
	bun.BaseModel `bun:"table:confirmation_of_payee_token"`

	ID        int64  `bun:",pk"`
	AccountID int64  `bun:",notnull"`
	Hash      string `bun:",unique,notnull"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UsedAt    time.Time `bun:",notnull"`
}
