package user

import (
	"github.com/sirupsen/logrus"
	"github.com/go-chi/chi"
	"net/http"
	"github.com/batazor/go-auth/models"
	"github.com/batazor/go-auth/db"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"github.com/batazor/go-auth/utils"
	"time"
	"errors"
)

var log = logrus.New()

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)
}

// Error handler
func Error(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	log.Error(err)

	err_str := `{
		"success": false,
		"error": [
			'` + err.Error()  + `'
		]
	}`

	w.Write([]byte(err_str))
	return
}

func CheckUniqueUser(w http.ResponseWriter, user models.User) bool {
	count, err := db.Session.DB("users").C(models.CollectionUser).Find(bson.M{"mail": user.Mail}).Count()
	if err != nil {
		Error(w, err)
		return true
	}

	if (count > 0) {
		w.WriteHeader(http.StatusBadRequest)
		Error(w, errors.New("need unique mail"))
		return true
	}

	return false
}

// Routes creates a REST router
func Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", List)
	r.Post("/", Create)
	r.Put("/{userId}", Update)
	r.Delete("/{userId}", Delete)

	return r
}

func List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := []models.User{}
	err := db.Session.DB("users").C(models.CollectionUser).Find(nil).Sort("-updated_on").All(&users)
	if err != nil {
		Error(w, err)
		return
	}

	res, err := json.Marshal(&users)
	if err != nil {
		Error(w, err)
	}

	w.Write(res)
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		Error(w, err)
		return
	}

	var user models.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		Error(w, err)
		return
	}

	is_err := CheckUniqueUser(w, user)
	if is_err { return }

	id := bson.NewObjectId()
	user.Id = id
	user.Password, _ = utils.HashPassword(user.Password)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err = db.Session.DB("users").C(models.CollectionUser).Insert(user)
	if err != nil {
		Error(w, err)
		return
	}

	output, err := json.Marshal(user)
	if err != nil {
		Error(w, err)
		return
	}

	w.Write(output)
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		Error(w, err)
		return
	}

	var user models.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		Error(w, err)
		return
	}

	var userId = chi.URLParam(r, "userId")
	if len(userId) != 24 {
		Error(w, errors.New("not correct user id"))
		return
	}

	user.Password, _ = utils.HashPassword(user.Password)
	user.UpdatedAt = time.Now()

	err = db.Session.DB("users").C(models.CollectionUser).UpdateId(bson.ObjectIdHex(userId), user)
	if err != nil {
		Error(w, err)
		return
	}

	output, err := json.Marshal(user)
	if err != nil {
		Error(w, err)
		return
	}

	w.Write(output)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userId = chi.URLParam(r, "userId")
	defer func() {
		var err = db.Session.DB("users").C(models.CollectionUser).RemoveId(bson.ObjectIdHex(userId))
		if err != nil {
			Error(w, err)
			return
		}
		recover()
	}()

	w.Write([]byte("{\"success\": true}"))
}