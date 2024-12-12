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

func bulkInsert(db *bun.DB) {
	ctx := context.Background()
	for _, query := range []*bun.CreateTableQuery{
		db.NewCreateTable().Model((*User)(nil)).IfNotExists(),
		db.NewCreateTable().Model((*AuthenticationInfo)(nil)).IfNotExists(),
	} {
		_, err := query.Exec(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	// insert data
	johnDoe := &User{
		Name:  "John Doe",
		Age:   30,
		Birth: Must(time.Parse("2006-01-02", "1990-01-01")),
		AuthenticationInfo: &AuthenticationInfo{
			Email:    "johndoe@example.com",
			Password: "password",
		},
	}

	// log.Print(db.NewInsert().Model(johnDoe).String())
	err := db.NewInsert().Model(johnDoe).Scan(ctx)
	if err != nil {
		log.Fatal(err)
	}
	johnDoe.AuthenticationInfo.UserID = johnDoe.ID
	err = db.NewInsert().Model(johnDoe.AuthenticationInfo).Scan(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db := openDB()
	defer db.Close()

	bulkInsert(db)

	ctx := context.Background()

	var oldPassword string
	err := db.NewSelect().
		Model(&AuthenticationInfo{}).
		Column("password").
		Where("email = ?", "johndoe@example.com").
		Scan(ctx, &oldPassword)
	if err != nil {
		log.Fatal(err)
	}
	if oldPassword != "password" {
		log.Fatal("wrong password")
	}

	// update password
	authorization := &AuthenticationInfo{Password: "new_password", UpdatedAt: time.Now()}
	_, err = db.NewUpdate().
		Model(authorization).
		OmitZero().
		Where("email = ?", "johndoe@example.com").
		Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// select data
	var users []*User
	err = db.NewSelect().Model(&User{}).Relation("AuthenticationInfo").Scan(ctx, &users)
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		log.Printf("%+v , %+v", user, user.AuthenticationInfo)
	}

}
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
