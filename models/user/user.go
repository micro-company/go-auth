package userModel

import (
	"errors"
	"time"

	"github.com/micro-company/go-auth/db"
	"gopkg.in/mgo.v2/bson"
)

const (
	// CollectionUser holds the name of the articles collection
	CollectionUser = "users"
)

func List() (error, []User) {
	var users []User
	err := db.Session.DB("auth").C(CollectionUser).Find(nil).Sort("-updated_on").All(&users)
	if err != nil {
		return err, users
	}

	return nil, users
}

func Find(user User) (error, []User) {
	var users []User
	err := db.Session.DB("auth").C(CollectionUser).Find(user).Sort("-updated_on").All(&users)
	if err != nil {
		return err, users
	}

	return nil, users
}

func FindOne(user User) (User, error) {
	err := db.Session.DB("auth").C(CollectionUser).Find(user).One(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func FindCount(user User) (int, error) {
	count, err := db.Session.DB("auth").C(CollectionUser).Find(user).Count()
	if err != nil {
		return count, err
	}

	return count, nil
}

func Add(user User) (error, User) {
	var checkUser User
	checkUser.Mail = user.Mail
	result, err := FindCount(checkUser)
	if result > 0 {
		return err, user
	}

	err = db.Session.DB("auth").C(CollectionUser).Insert(user)
	if err != nil {
		return errors.New(`{"mail":"need uniq mail"}`), user
	}

	return nil, user
}

func Update(user User) (error, User) {
	UpdatedAt := time.Now()
	user.UpdatedAt = &UpdatedAt
	user.Mail = nil // prohibit changing address

	err := db.Session.DB("auth").C(CollectionUser).UpdateId(user.Id, bson.M{"$set": user})
	if err != nil {
		return err, user
	}

	return nil, user
}

func Delete(userId string) error {
	err := db.Session.DB("auth").C(CollectionUser).RemoveId(bson.ObjectIdHex(userId))
	if err != nil {
		return err
	}

	return nil
}
