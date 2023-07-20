package models

type UserLogin struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
