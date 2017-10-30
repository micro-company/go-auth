package user

import (
	"errors"
	"github.com/batazor/go-auth/db"
	"github.com/batazor/go-auth/utils"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"github.com/batazor/go-auth/models/user"
)

func CheckUniqueUser(w http.ResponseWriter, user userModel.User) bool {
	count, err := db.Session.DB("users").C(userModel.CollectionUser).Find(bson.M{"mail": user.Mail}).Count()
	if err != nil {
		utils.Error(w, err)
		return true
	}

	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)
		utils.Error(w, errors.New("need unique mail"))
		return true
	}

	return false
}
