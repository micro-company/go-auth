package userModel

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// User model
type User struct {
	Id            *bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Email         *string        `json:"email" bson:"email,omitempty"`
	VerifiedEmail bool           `json:"veriefied_email"`
	Locale        string         `json:"locale" bson:"locale"`
	Password      string         `json:"password" bson:"password,omitempty"`
	PasswordRetry string         `json:"retryPassword" bson:"-"`
	RecoveryToken string         `json:"recoveryToken" bson:"-"`
	Gender        string         `json:"gender" bson:"-"`
	Profiles      []Profile      `json:"profiles" bson:"profiles"`
	CreatedAt     *time.Time     `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt     *time.Time     `json:"updated_at" bson:"updated_at,omitempty"`
}

type Profile struct {
	Type    string `json:"type"`
	Id      string `json:"id"`
	Link    string `json:"link"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}
