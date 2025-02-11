package entity

import (
	"oauth-server/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	USER_TABLE_NAME = "users"
)

type User struct {
	gorm.Model
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email       *string   `json:"email" gorm:"type:varchar(100);"`
	PhoneNumber *string   `json:"phone_number" gorm:"type:varchar(20);"`
	Password    string    `json:"password" gorm:"type:text;not null"`
	IsActive    bool      `json:"is_active" gorm:"default:true;type:bool;not null"`
	CreatedAt   int64     `json:"created_at" gorm:"type:integer;not null"`
	UpdatedAt   int64     `json:"updated_at" gorm:"type:integer;not null"`
}

func NewUser() *User {
	return &User{
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (u *User) TableName() string {
	return USER_TABLE_NAME
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" {
		hash, err := utils.HashPassword(u.Password)
		if err != nil {
			return err
		}

		u.Password = hash
	}

	return nil
}

func (u *User) GetUserID() string {
	return u.ID.String()
}

func (u *User) GetPhoneNumber() string {
	return *u.PhoneNumber
}

func (u *User) GetEmail() string {
	return *u.Email
}
