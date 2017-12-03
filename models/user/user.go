package userModel

import (
	"time"

	"github.com/micro-company/go-auth/db"
	"gopkg.in/mgo.v2/bson"
)

const (
	// CollectionUser holds the name of the articles collection
	CollectionUser = "users"
)

// User model
type User struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Mail      string        `json:"mail" bson:"mail,omitempty"`
	Password  string        `json:"password" bson:"password,omitempty"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at,omitempty"`
}

func List() (error, []User) {
	var users []User
	err := db.Session.DB("users").C(CollectionUser).Find(nil).Sort("-updated_on").All(&users)
	if err != nil {
		return err, users
	}

	return nil, users
}

func Find(user User) (error, []User) {
	var users []User
	err := db.Session.DB("users").C(CollectionUser).Find(user).Sort("-updated_on").All(&users)
	if err != nil {
		return err, users
	}

	return nil, users
}

func FindOne(user User) (error, User) {
	err := db.Session.DB("users").C(CollectionUser).Find(&user).One(&user)
	if err != nil {
		return err, user
	}

	return nil, user
}

func Add(user User) (error, User) {
	err := db.Session.DB("users").C(CollectionUser).Insert(user)
	if err != nil {
		return err, user
	}

	return nil, user
}

func Update(user User) (error, User) {
	err := db.Session.DB("users").C(CollectionUser).UpdateId(user.Id, user)
	if err != nil {
		return err, user
	}

	return nil, user
}

func Delete(userId string) error {
	err := db.Session.DB("users").C(CollectionUser).RemoveId(bson.ObjectIdHex(userId))
	if err != nil {
		return err
	}

	return nil
}
