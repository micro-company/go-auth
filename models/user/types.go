package userModel

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// User model
type User struct {
	Id            *bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Mail          *string        `json:"mail" bson:"mail,omitempty"`
	MailVerified  bool           `json:"mail_verified"`
	Password      string         `json:"password" bson:"password,omitempty"`
	PasswordRetry string         `json:"retryPassword" bson:"-"`
	RecoveryToken string         `json:"recoveryToken" bson:"-"`
	CreatedAt     *time.Time     `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt     *time.Time     `json:"updated_at" bson:"updated_at,omitempty"`
}
