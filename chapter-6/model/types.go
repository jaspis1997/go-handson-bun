package model

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel      `bun:"table:users"`
	ID                 int64               `bun:"id,pk,autoincrement"`
	Name               string              `bun:"name,notnull"`
	Age                int                 `bun:"age,notnull"`
	Birth              time.Time           `bun:"birth,nullzero"`
	AuthenticationInfo *AuthenticationInfo `bun:"rel:has-one,join:id=user_id"`
}

type AuthenticationInfo struct {
	bun.BaseModel `bun:"table:authentication_infos"`
	UserID        int64     `bun:"user_id,pk,notnull"`
	Email         string    `bun:"email,unique,notnull"`
	Password      string    `bun:"password,notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,default:current_timestamp"`
}
