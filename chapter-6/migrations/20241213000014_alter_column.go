package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	up := func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewRaw("alter table users add column deleted bool").Exec(ctx)
		return err
	}
	down := func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewRaw("alter table users drop column deleted").Exec(ctx)
		return err
	}
	MigrationGroup.MustRegister(up, down)
}
