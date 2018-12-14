package userModel

import (
	"errors"
	"time"

	"github.com/micro-company/go-auth/db/mongodb"
	"github.com/micro-company/go-auth/utils/crypto"
	"github.com/mongodb/mongo-go-driver/bson"
)

const (
	// CollectionUser holds the name of the articles collection
	CollectionUser = "users"
)

func List() (error, []User) {
	var users []User
	//err := mongodb.Session.Database("auth").Collection(CollectionUser).Find(nil).Sort("-updated_on").All(&users)
	//cursor, err := mongodb.Session.Database("auth").Collection(CollectionUser).Find(nil, users)
	//if err != nil {
	//	return err, users
	//}

	return nil, users
}

func Find(user User) (*[]User, error) {
	var users *[]User
	res := mongodb.Session.Database("auth").Collection(CollectionUser).FindOne(nil, user)
	res.Decode(users)
	if res.Err() != nil {
		return nil, res.Err()
	}

	return users, nil
}

func FindOne(user User) (*User, error) {
	res := mongodb.Session.Database("auth").Collection(CollectionUser).FindOne(nil, user)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var us *User
	res.Decode(us)
	return us, nil
}

func FindCount(user User) (int64, error) {
	count, err := mongodb.Session.Database("auth").Collection(CollectionUser).Count(nil, user)
	if err != nil {
		return count, err
	}

	return count, nil
}

func Add(user User) (error, User) {
	var checkUser User
	checkUser.Email = user.Email
	result, err := FindCount(checkUser)
	if result > 0 {
		return err, user
	}

	user.Password, _ = crypto.HashPassword(user.Password)
	time := time.Now()
	user.UpdatedAt = &time

	res, err := mongodb.Session.Database("auth").Collection(CollectionUser).InsertOne(nil, user)
	if err != nil {
		return errors.New(`{"mail":"need uniq mail"}`), user
	}

	user.Id = res.InsertedID.(string)

	return nil, user
}

func Update(user *User) (*User, error) {
	UpdatedAt := time.Now()
	user.UpdatedAt = &UpdatedAt
	user.Email = nil // prohibit changing address

	filter := bson.D{{"_id", user.Id}}
	_, err := mongodb.Session.Database("auth").Collection(CollectionUser).UpdateOne(nil, filter, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Delete(userId string) (int64, error) {
	filter := bson.D{{"_id", userId}}
	res, err := mongodb.Session.Database("auth").Collection(CollectionUser).DeleteOne(nil, filter)

	return res.DeletedCount, err
}
