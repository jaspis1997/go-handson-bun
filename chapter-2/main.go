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
}

func openDB() *bun.DB {
	return bun.NewDB(
		Must(sql.Open(sqliteshim.ShimName, "file:test.db")),
		sqlitedialect.New(),
	)
}

func bulkInsert(db *bun.DB) {
	ctx := context.Background()
	_, err := db.NewCreateTable().Model((*User)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// insert data
	johnDoe := &User{
		Name:  "John Doe",
		Age:   30,
		Birth: Must(time.Parse("2006-01-02", "1990-01-01")),
	}

	// log.Print(db.NewInsert().Model(johnDoe).String())
	err = db.NewInsert().Model(johnDoe).Scan(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db := openDB()
	defer db.Close()

	bulkInsert(db)

	ctx := context.Background()
	// select data
	var names []string
	err := db.NewSelect().Model(&User{}).Column("name").Scan(ctx, &names)
	if err != nil {
		log.Fatal(err)
	}
	for _, name := range names {
		log.Printf("%s", name)
	}

	var users []*User
	err = db.NewSelect().Model(&User{}).Column("name", "age").Scan(ctx, &users)
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		log.Printf("%v", user)
	}

}
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
