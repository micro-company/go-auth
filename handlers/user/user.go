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
)

var log = logrus.New()

type error interface {
	Error() string
}

// Error catch
func Error(w http.ResponseWriter, err error) {
	log.Error(err)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"success\": false}"))
	return
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
		log.Fatal(err)
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

	// Check unique mail
	count, err := db.Session.DB("users").C(models.CollectionUser).Find(bson.M{"mail": user.Mail}).Count()
	if err != nil {
		Error(w, err)
		return
	}

	if (count > 0) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"success": false,
			"error": [
				"need unique mail"
			]
		}`))
		return
	}

	id := bson.NewObjectId()
	user.Id = id
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
	var err = db.Session.DB("users").C(models.CollectionUser).RemoveId(bson.ObjectIdHex(userId))
	if err != nil {
		Error(w, err)
		return
	}

	w.Write([]byte("{\"success\": true}"))
}