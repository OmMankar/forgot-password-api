package models

type User struct {
	EmailId  string `bson:"emailId" json:"emailId" validation:"required"`
	Password string `bson:"password" json:"password" validation:"required"`
}
