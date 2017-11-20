package sessionModel

import (
	"time"

	"github.com/batazor/go-auth/db"
)

// Session model
type Session struct {
	token  string
	status bool
}

func Add(token string, status bool, ttl time.Duration) error {
	err := db.Redis.Set(token, status, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func Delete(token string) error {
	err := db.Redis.Del(token).Err()
	if err != nil {
		return err
	}

	return nil
}
