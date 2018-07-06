package recaptcha

import (
	"encoding/json"
	"errors"
	"github.com/micro-company/go-auth/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	RECAPTCHA_PRIVATE_KEY string
)

const recaptchaServerName = "https://www.google.com/recaptcha/api/siteverify"

func init() {
	// Get configuration
	RECAPTCHA_PRIVATE_KEY = utils.Getenv("RECAPTCHA_PRIVATE_KEY", "localhost/auth")
}

func VerifyCaptcha(captchaResponse []byte) error {
	var captcha Captcha
	err := json.Unmarshal(captchaResponse, &captcha)
	if err != nil {
		return err
	}

	r, err := http.PostForm(recaptchaServerName,
		url.Values{"secret": {RECAPTCHA_PRIVATE_KEY}, "response": {captcha.Captcha}})
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return errors.New(`"` + err.Error() + `"`)
	}

	var recaptcha RecaptchaResponse
	err = json.Unmarshal(b, &recaptcha)
	if err != nil {
		return errors.New(`"` + err.Error() + `"`)
	}

	if recaptcha.Success {
		return nil
	}

	return errors.New(`"` + strings.Join(recaptcha.ErrorCodes[:], ",") + `"`)
}
