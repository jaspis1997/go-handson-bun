package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/migrate"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            int       `bun:"id,pk,autoincrement"`
	Name          string    `bun:"name,notnull"`
	Age           int       `bun:"age,notnull"`
	Birth         time.Time `bun:"birth,nullzero"`

	CreatedAt time.Time `bun:"created_at,nullzero,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,default:current_timestamp"`

	AuthenticationInfo *AuthenticationInfo `bun:"rel:has-one,join:id=user_id"`
}

type AuthenticationInfo struct {
	bun.BaseModel `bun:"table:authentication_infos"`
	UserID        int       `bun:"user_id,pk,notnull"`
	Email         string    `bun:"email,unique,notnull"`
	Password      string    `bun:"password,notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,default:current_timestamp"`
}

func openDB() *bun.DB {
	return bun.NewDB(
		Must(sql.Open(sqliteshim.ShimName, "file:test.db")),
		sqlitedialect.New(),
	)
}

var migrations = migrate.NewMigrations()

func main() {
	db := openDB()
	ctx := context.Background()

	migrator := migrate.NewMigrator(db, migrations)
	if err := migrator.Init(ctx); err != nil {
		log.Fatal(err)
	}

	migrator.Lock(ctx)
	defer migrator.Unlock(ctx)

	group, err := migrator.Migrate(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if group.IsZero() {
		log.Print("no migrations found")
		return
	}
	log.Printf("migrated %s", group)
}

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
