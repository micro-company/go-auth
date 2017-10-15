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

// Routes creates a REST router
func Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", List)
	r.Post("/", Create)

	return r
}

func List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := []models.User{}
	err := db.Session.DB("users").C(models.CollectionUser).Find(nil).Sort("-updated_on").All(&users)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"success\": false}"))
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
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"success\": false}"))
		return
	}

	var book models.User
	err = json.Unmarshal(b, &book)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"success\": false}"))
		return
	}

	id := bson.NewObjectId()
	book.Id = id
	err = db.Session.DB("books").C(models.CollectionUser).Insert(book)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"success\": false}"))
		return
	}

	output, err := json.Marshal(book)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"success\": false}"))
		return
	}

	w.Write(output)
}