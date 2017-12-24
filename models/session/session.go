package sessionModel

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/micro-company/go-auth/db"
	"github.com/sirupsen/logrus"
)

const (
	PRIVATE_KEY = "cert/private_key.pem"
	PUBLIC_KEY  = "cert/public_key.pub"
)

var (
	log = logrus.New()

	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

// Session model
type Session struct {
	token  string
	status bool
}

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)

	// JWT =====================================================================
	signBytes, err := ioutil.ReadFile(PRIVATE_KEY)
	if err != nil {
		log.Fatal(err)
		return
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
		return
	}

	verifyBytes, err := ioutil.ReadFile(PUBLIC_KEY)
	if err != nil {
		log.Fatal(err)
		return
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func NewAccessToken(timeDuration int64) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = timeDuration
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewRefreshToken(timeDuration time.Duration) (string, error) {
	refreshToken, _ := uuid.NewUUID()
	err := db.Redis.Set(refreshToken.String(), "true", timeDuration).Err()
	if err != nil {
		return "", err
	}

	return refreshToken.String(), nil
}

func NewRecoveryLink(value string) (string, error) {
	TTL := time.Hour * 1
	refreshToken, _ := uuid.NewUUID()
	log.Info("NewRecoveryLink", value)
	log.Info("refreshToken", refreshToken.String())
	err := db.Redis.Set(refreshToken.String(), value, TTL).Err()
	if err != nil {
		return "", err
	}

	return refreshToken.String(), nil
}

func Delete(token string) error {
	err := db.Redis.Del(token).Err()
	if err != nil {
		return err
	}

	return nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return verifyKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, err
}

func CheckRefreshToken(token string) (bool, error) {
	value := db.Redis.Get(token)
	if value.Err() != nil {
		return false, value.Err()
	}

	status, err := value.Result()
	if err != nil && status != "true" {
		return false, err
	}

	return true, nil
}

func GetValueByKey(token string) (string, error) {
	value := db.Redis.Get(token)
	if value.Err() != nil {
		return "", value.Err()
	}

	status, err := value.Result()
	if err != nil && status != "true" {
		return "", err
	}

	return status, nil
}
