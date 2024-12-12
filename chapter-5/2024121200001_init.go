package main

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	up := func(ctx context.Context, db *bun.DB) error {
		for _, query := range []*bun.CreateTableQuery{
			db.NewCreateTable().Model((*User)(nil)).IfNotExists(),
			db.NewCreateTable().Model((*AuthenticationInfo)(nil)).IfNotExists(),
		} {
			_, err := query.Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}
	down := func(ctx context.Context, db *bun.DB) error {
		for _, query := range []*bun.DropTableQuery{
			db.NewDropTable().Model((*User)(nil)),
			db.NewDropTable().Model((*AuthenticationInfo)(nil)),
		} {
			_, err := query.Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}
	migrations.MustRegister(up, down)
}
