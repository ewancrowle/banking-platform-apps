package device

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Device struct {
	bun.BaseModel `bun:"table:devices"`

	ID int64 `bun:",pk"`

	AccountID *int64 `bun:"account_id,notnull"`
	//Account   *account.Account `bun:"rel:belongs-to,join:account_id=id"`

	IPAddress string `bun:"ip_address,notnull"`
	UserAgent string `bun:"user_agent,notnull"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

func (d *Device) Insert(ctx context.Context, db *bun.DB) error {
	_, err := db.NewInsert().Model(d).Exec(ctx)
	return err
}

func Select(ctx context.Context, db *bun.DB, id int64) (*Device, error) {
	d := new(Device)
	err := db.NewSelect().Model(d).Where("id = ?", id).Scan(ctx)
	return d, err
}

func (d *Device) Delete(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDelete().Model(d).WherePK().Exec(ctx)
	return err
}
