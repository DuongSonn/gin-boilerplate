package model

type LoginRequest struct {
	Password    string  `json:"password" validate:"required"`
	Email       *string `json:"email" validate:"omitempty,email"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,phone_number"`
}
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	Password        string  `json:"password" validate:"required"`
	ConfirmPassword string  `json:"confirm_password" validate:"required,eqfield=Password"`
	Email           *string `json:"email" validate:"omitempty,email"`
	PhoneNumber     *string `json:"phone_number" validate:"omitempty,phone_number"`
}
type RegisterResponse struct {
}

type LogoutRequest struct {
}
type LogoutResponse struct {
}
