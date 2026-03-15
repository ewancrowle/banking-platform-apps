package migrations

import (
	"account-svc/pkg/model/account"
	"context"
	"fmt"
	"oauth-svc/pkg/model/device"
	"oauth-svc/pkg/model/token"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")

		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			_, err := tx.NewCreateTable().
				Model((*account.Account)(nil)).
				IfNotExists().
				Exec(ctx)
			if err != nil {
				return err
			}

			_, err = tx.NewCreateTable().
				Model((*device.Device)(nil)).
				IfNotExists().
				WithForeignKeys().
				Exec(ctx)
			if err != nil {
				return err
			}

			_, err = tx.NewCreateTable().
				Model((*token.AccessToken)(nil)).
				IfNotExists().
				WithForeignKeys().
				Exec(ctx)
			if err != nil {
				return err
			}

			_, err = tx.NewCreateTable().
				Model((*token.RefreshToken)(nil)).
				IfNotExists().
				WithForeignKeys().
				Exec(ctx)
			if err != nil {
				return err
			}

			return err
		})
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")

		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			_, err := tx.NewDropTable().
				Model((*token.RefreshToken)(nil)).
				IfExists().
				Exec(ctx)
			if err != nil {
				return err
			}

			_, err = tx.NewDropTable().
				Model((*token.AccessToken)(nil)).
				IfExists().
				Exec(ctx)
			if err != nil {
				return err
			}

			_, err = tx.NewDropTable().
				Model((*device.Device)(nil)).
				IfExists().
				Exec(ctx)
			if err != nil {
				return err
			}

			_, err = tx.NewDropTable().
				Model((*account.Account)(nil)).
				IfExists().
				Cascade().
				Exec(ctx)
			if err != nil {
				return err
			}

			return err
		})
	})
}
