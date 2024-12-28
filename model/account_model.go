package model

type CreateAccountRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

type ResponseAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
