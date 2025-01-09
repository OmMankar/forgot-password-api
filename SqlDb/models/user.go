package models

type User struct {
	EmailId  string `json:"emailId" validation:"required"`
	PassWord string `json:"password" validation:"required"`
}
