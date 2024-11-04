package model

import "github.com/google/uuid"

type CreateAccountRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

type ResponseAccount struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}
