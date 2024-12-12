package migrations

import (
	"context"
	"handson/chapter-6/model"

	"github.com/uptrace/bun"
)

func init() {
	up := func(ctx context.Context, db *bun.DB) error {
		for _, query := range []*bun.CreateTableQuery{
			db.NewCreateTable().Model((*model.User)(nil)).IfNotExists(),
			db.NewCreateTable().Model((*model.AuthenticationInfo)(nil)).IfNotExists(),
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
			db.NewDropTable().Model((*model.User)(nil)),
			db.NewDropTable().Model((*model.AuthenticationInfo)(nil)),
		} {
			_, err := query.Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}
	MigrationGroup.MustRegister(up, down)
}
