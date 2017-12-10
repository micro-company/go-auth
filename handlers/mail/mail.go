package mail

import (
	"bytes"
	"html/template"
	"net/smtp"

	"github.com/sirupsen/logrus"

	"github.com/micro-company/go-auth/utils"
)

var (
	log           = logrus.New()
	auth          smtp.Auth
	emailUser     EmailUser
	SMTP_USERNAME = utils.Getenv("SMTP_USERNAME", "example@gmail.com")
	SMTP_PASSWORD = utils.Getenv("SMTP_PASSWORD", "secretPass")
	SMTP_SERVER   = utils.Getenv("SMTP_SERVER", "smtp.gmail.com")
	SMTP_PORT     = utils.Getenv("SMTP_PORT", "587")
	emailTemplate string
)

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)
}

func Recovery(data RecoveryData) error {
	emailUser := EmailUser{SMTP_USERNAME, SMTP_PASSWORD, SMTP_SERVER, SMTP_PORT}
	auth = smtp.PlainAuth("",
		emailUser.Username,
		emailUser.Password,
		emailUser.Server,
	)

	emailTemplate, err := ParseTemplate("recovery.html", data)
	if err != nil {
		return err
	}

	subject := "Subject: Recovery pass!"

	err = Send(emailUser, subject, emailTemplate)
	if err != nil {
		return err
	}

	return nil
}

func Send(emailUser EmailUser, subject, emailTemplate string) error {
	addr := emailUser.Server + ":" + emailUser.Port
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(subject + "\n" + mime + "\n" + emailTemplate)

	err := smtp.SendMail(
		addr,
		auth,
		emailUser.Username,
		[]string{"batazor111@gmail.com"},
		msg,
	)
	if err != nil {
		return err
	}
	return nil
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles("handlers/mail/template/" + templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
