package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	// CollectionUser holds the name of the articles collection
	CollectionUser = "users"
)

// User model
type User struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Mail      string        `json:"mail" bson:"mail"`
	Password  string        `json:"password" bson:"password"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}
