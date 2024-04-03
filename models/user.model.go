package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id           primitive.ObjectID   `json:"_id" bson:"_id"`
	Firstname    string               `json:"firstName" bson:"firstName"`
	Lastname     string               `json:"lastName" bson:"lastName"`
	Email        string               `json:"email" bson:"email"`
	Password     string               `json:"password" bson:"password"`
	MySpaces     []primitive.ObjectID `json:"my_spaces,omitempty" bson:"my_spaces,omitempty"`
	SharedSpaces []primitive.ObjectID `json:"shared_spaces" bson:"shared_spaces"`
}

type Role struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Role        string             `json:"role" bson:"role"`
	Permissions []string           `json:"permissions" bson:"permissions"`
}

type SpaceUsers struct {
	UserId primitive.ObjectID `json:"_id" bson:"_id"`
	Role   string             `json:"role" bson:"role"`
}

type Space struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Users     []SpaceUsers       `json:"users" bson:"users"`
	Name      string             `json:"name" bson:"name"`
	Admin     primitive.ObjectID `json:"admin" bson:"admin"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
}

type Comment struct {
	Id      primitive.ObjectID `json:"_id" bson:"_id"`
	SentBy  primitive.ObjectID `json:"sentBy" bson:"sentBy"`
	SentTo  primitive.ObjectID `json:"sentTo" bson:"sentTo"`
	Content interface{}        `json:"content" bson:"content"`
	SentAt  primitive.DateTime `json:"sentAt" bson:"sentAt"`
	SpaceId primitive.ObjectID `json:"spaceId" bson:"spaceId"`
}

type Login struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type LoginSuccess struct {
	Id           primitive.ObjectID   `json:"id"`
	Firstname    string               `json:"firstName" bson:"firstName"`
	Lastname     string               `json:"lastName" bson:"lastName"`
	Email        string               `json:"email"`
	UUID         string               `json:"userUUID"`
	MySpaces     []primitive.ObjectID `json:"my_spaces"`
	SharedSpaces []primitive.ObjectID `json:"shared_spaces"`
}

type UserSpaces struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id"`
	MySpaces     []Space            `json:"mySpaces" bson:"mySpaces"`
	SharedSpaces []Space            `json:"sharedSpaces" bson:"sharedSpaces"`
}

type Verify struct {
	Token string `json:"token"`
}

type PromptPayload struct {
	Content string `json:"content" bson:"content"`
}

type Message struct {
	Prompt string `json:"prompt" bson:"prompt"`
}
