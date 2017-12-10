package mail

type EmailUser struct {
	Username string
	Password string
	Server   string
	Port     string
}

type RecoveryData struct {
	Mail string
	Url  string
}
