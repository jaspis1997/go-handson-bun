package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	up := func(ctx context.Context, db *bun.DB) error {
		return nil
	}
	down := func(ctx context.Context, db *bun.DB) error {
		return nil
	}
	migrations.MustRegister(up, down)
}
