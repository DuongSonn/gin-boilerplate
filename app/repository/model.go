package repository

import "github.com/google/uuid"

type FindUserByFilter struct {
	Email        *string
	PhoneNumber  *string
	ID           *uuid.UUID
	IDs          []uuid.UUID
	Emails       []string
	PhoneNumbers []string
	Limit        *int
	Offset       *int
	IsActive     *bool
}

type FindOAuthByFilter struct {
	UserID   *uuid.UUID
	Token    *string
	PlatForm *string
}
