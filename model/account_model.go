package model

type CreateAccountRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

type ResponseAccount struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
