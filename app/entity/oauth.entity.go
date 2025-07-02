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

type OAuthStatus string

const (
	OAuthStatusActive   OAuthStatus = "active"
	OAuthStatusInactive OAuthStatus = "inactive"
	OAuthStatusBlocked  OAuthStatus = "blocked"
)

type OAuth struct {
	ID        uuid.UUID   `json:"id" gorm:"primaryKey;type:uuid"`
	UserID    uuid.UUID   `json:"user_id" gorm:"type:uuid;not null"`
	IP        string      `json:"ip" gorm:"type:text;not null"`
	Platform  string      `json:"platform" gorm:"type:varchar(10);not null"`
	Token     string      `json:"token" gorm:"type:text;not null"`
	Status    OAuthStatus `json:"status" gorm:"varchar(10);not null"`
	ExpireAt  int64       `json:"expire_at" gorm:"type:integer;not null"`
	CreatedAt int64       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64       `json:"updated_at" gorm:"autoUpdateTime:milli"`
	LoginAt   int64       `json:"login_at" gorm:"type:integer"`
}

func NewOAuth() *OAuth {
	return &OAuth{
		ID:        uuid.New(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (OAuth) TableName() string {
	return OAUTH_TABLE_NAME
}
