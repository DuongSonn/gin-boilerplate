package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	OAUTH_TABLE_NAME = "oauth"
)

const (
	OAuthPlatformMobile = "mobile"
	OAuthPlatformWeb    = "web"
)

const (
	OAuthStatusActive   = "active"
	OAuthStatusInactive = "inactive"
	OAuthStatusBlocked  = "blocked"
)

type Oauth struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	IP        string    `json:"ip" gorm:"type:text;not null"`
	Platform  string    `json:"platform" gorm:"type:varchar(10);not null"`
	Token     string    `json:"token" gorm:"type:text;not null"`
	Status    string    `json:"status" gorm:"varchar(10);not null"`
	ExpireAt  int64     `json:"expire_at" gorm:"type:integer;not null"`
	CreatedAt int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`
	LoginAt   int64     `json:"login_at" gorm:"type:integer"`
}

func NewOAuth() *Oauth {
	return &Oauth{
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (Oauth) TableName() string {
	return OAUTH_TABLE_NAME
}
