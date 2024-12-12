package main

import (
	"database/sql"
	"fmt"
	"handson/chapter-6/migrations"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

var rootCommand = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Root command")
	},
}

const Version = "v0.0.1-builtin"

func init() {
	rootCommand.AddCommand(&cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	})
	migrateCommand := &cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			log.Print("Migrate command")
			db := openDB()
			defer db.Close()
			migrations.Migrate(db)
		},
	}
	rootCommand.AddCommand(migrateCommand)
	migrateCommand.AddCommand(&cobra.Command{
		Use:     "create",
		Example: "migrate create ...$comment ",

		Run: func(cmd *cobra.Command, args []string) {
			log.Print("Create command", args)
			for _, arg := range args {
				migrations.CreateFile(arg)
			}
		},
	})
	migrateCommand.AddCommand(&cobra.Command{
		Use: "status",
		Run: func(cmd *cobra.Command, args []string) {
			db := openDB()
			defer db.Close()
			migrations.Status(db)
		},
	})
	migrateCommand.AddCommand(&cobra.Command{
		Use: "rollback",
		Run: func(cmd *cobra.Command, args []string) {
			log.Print("Rollback command")
			db := openDB()
			defer db.Close()
			migrations.Rollback(db)
		},
	})
}

func main() {
	rootCommand.Execute()
}

func openDB() *bun.DB {
	return bun.NewDB(
		Must(sql.Open(sqliteshim.ShimName, "file:test.db")),
		sqlitedialect.New(),
	)
}

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
