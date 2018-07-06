package user

import (
	"errors"
	"net/http"

	"github.com/micro-company/go-auth/db"
	"github.com/micro-company/go-auth/models/user"
	"github.com/micro-company/go-auth/utils"
	"gopkg.in/mgo.v2/bson"
)

func CheckUniqueUser(w http.ResponseWriter, user userModel.User) bool {
	count, err := db.Session.DB("auth").C(userModel.CollectionUser).Find(bson.M{"mail": user.Email}).Count()
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return true
	}

	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)
		utils.Error(w, errors.New(`{"mail": "need unique mail"}`))
		return true
	}

	return false
}
