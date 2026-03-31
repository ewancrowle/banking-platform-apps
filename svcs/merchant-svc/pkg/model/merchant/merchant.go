package merchant

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Merchant struct {
	bun.BaseModel `bun:"table:merchants"`

	ID int64 `bun:",pk"`

	Descriptor      string `bun:",notnull"`
	ShortDescriptor string `bun:",notnull"`

	// See https://www.iso.org/standard/79450.html
	// Also https://www.mastercard.com/content/dam/mccom/shared/business/support/rules-pdfs/mastercard-quick-reference-booklet-merchant.pdf
	MCC string `bun:",notnull"`

	Line1       string `bun:",notnull"`
	Line2       string
	Town        string `bun:",notnull"`
	Postcode    string `bun:",notnull"`
	CountryCode string `bun:",notnull"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func (m *Merchant) Insert(ctx context.Context, db *bun.DB) error {
	_, err := db.NewInsert().Model(m).Exec(ctx)
	return err
}

func Select(ctx context.Context, db *bun.DB, id int64) (*Merchant, error) {
	m := new(Merchant)
	err := db.NewSelect().Model(m).Where("id = ?", id).Scan(ctx)
	return m, err
}
