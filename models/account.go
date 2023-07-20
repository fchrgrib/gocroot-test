package models

type UserLogin struct {
	UserName string `json:"user_name,omitempty" bson:"user_name,omitempty"`
	Password string `json:"user_pass,omitempty" bson:"user_pass,omitempty"`
}
