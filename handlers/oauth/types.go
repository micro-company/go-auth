package oauth

type Callback struct {
	Code string `json:"code"`
}

type UserGoogle struct {
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	Gender        string `json:"gender"`
	GivenName     string `json:"given_name"`
	Id            string `json:"id"`
	Link          string `json:"link"`
	Locale        string `json:"locale"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"veriefied_email"`
}
