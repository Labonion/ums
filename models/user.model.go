package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Firstname string             `json:"firstName" bson:"firstName"`
	Lastname  string             `json:"lastName" bson:"lastName"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
}

type Login struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type LoginSuccess struct {
	Id        primitive.ObjectID `json:"id"`
	Firstname string             `json:"firstName" bson:"firstName"`
	Lastname  string             `json:"lastName" bson:"lastName"`
	Email     string             `json:"email"`
	UUID      string             `json:"userUUID"`
}

type Verify struct {
	Token string `json:"token"`
}
