package user

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/micro-company/go-auth/middleware"

	"github.com/micro-company/go-auth/utils/crypto"

	"github.com/go-chi/chi"
	"github.com/micro-company/go-auth/models/user"
	"github.com/micro-company/go-auth/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

var log = logrus.New()

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)
}

// Routes creates a REST router
func Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.CheckAuth)

	r.Get("/", List)
	r.Post("/", Create)
	r.Patch("/{userId}", Update)
	r.Delete("/{userId}", Delete)

	return r
}

func List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parent := opentracing.GlobalTracer().StartSpan("GET /users")
	defer parent.Finish()

	err, users := userModel.List()
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	res, err := json.Marshal(&users)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	w.Write(res)
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parent := opentracing.GlobalTracer().StartSpan("POST /users")
	defer parent.Finish()

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	var user userModel.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	is_err := CheckUniqueUser(w, user)
	if is_err {
		return
	}

	id := bson.NewObjectId()
	user.Id = &id
	user.Password, _ = crypto.HashPassword(user.Password)
	time := time.Now()
	user.CreatedAt = &time
	user.UpdatedAt = &time

	err, user = userModel.Add(user)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	output, err := json.Marshal(user)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(output)
	return
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parent := opentracing.GlobalTracer().StartSpan("PUT /users")
	defer parent.Finish()

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	var user userModel.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	var userId = chi.URLParam(r, "userId")
	if len(userId) != 24 {
		utils.Error(w, errors.New("not correct user id"))
		return
	}

	*user.Id = bson.ObjectIdHex(userId)
	user.Password, _ = crypto.HashPassword(user.Password)

	err, user = userModel.Update(user)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	output, err := json.Marshal(user)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	w.Write(output)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parent := opentracing.GlobalTracer().StartSpan("DELETE /users")
	defer parent.Finish()

	var userId = chi.URLParam(r, "userId")
	err := userModel.Delete(userId)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	w.Write([]byte(`{"id": "` + userId + `"}`))
}
