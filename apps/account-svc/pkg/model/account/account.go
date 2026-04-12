package account

import (
	"time"

	"github.com/uptrace/bun"
)

type Account struct {
	bun.BaseModel `bun:"table:accounts"`

	ID        int64  `bun:",pk"`
	AccountNo string `bun:",unique,notnull"`

	FirstName   string `bun:",notnull"`
	MiddleNames string
	LastName    string `bun:",notnull"`

	Email        string `bun:",unique,notnull"`
	PasswordHash string `bun:",notnull"`

	Line1    string `bun:",notnull"`
	Line2    string
	Town     string `bun:",notnull"`
	Postcode string `bun:",notnull"`

	IsFrozen     bool   `bun:",notnull"`
	CurrencyCode string `bun:",notnull"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}
