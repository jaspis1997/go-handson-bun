package migrations

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	_ "embed"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

var MigrationGroup = migrate.NewMigrations()

func initialize(migrator *migrate.Migrator) {
	if err := migrator.Init(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func Migrate(db *bun.DB) {
	ctx := context.Background()
	migrator := migrate.NewMigrator(db, MigrationGroup)
	initialize(migrator)
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

func Status(db *bun.DB) {
	migrator := migrate.NewMigrator(db, MigrationGroup)
	initialize(migrator)
	ms, err := migrator.MigrationsWithStatus(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("migrations: %s\n", ms)
	fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
	fmt.Printf("last migration group: %s\n", ms.LastGroup())
}

func Rollback(db *bun.DB) {
	migrator := migrate.NewMigrator(db, MigrationGroup)
	initialize(migrator)
	group, err := migrator.Rollback(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if group.ID == 0 {
		fmt.Printf("there are no groups to roll back\n")
		return
	}

	fmt.Printf("rolled back %s\n", group)
}

//go:embed __template.go
var migrationFileTemplate []byte

func CreateFile(comment string) {
	name := func() string {
		now := time.Now()
		timestamp := now.Format("20060102150405")
		name := strings.Join([]string{timestamp, comment}, "_") + ".go"
		return path.Join("migrations", name)
	}

	f, err := os.Create(name())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write(migrationFileTemplate)
	if err != nil {
		log.Fatal(err)
	}
}
