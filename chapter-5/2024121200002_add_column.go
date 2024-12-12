package main

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	up := func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewRaw("alter table users add column Deleted bool").Exec(ctx)
		return err
	}
	down := func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewRaw("alter table users drop column deleted").Exec(ctx)
		return err
	}
	migrations.MustRegister(up, down)
}
