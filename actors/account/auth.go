package account

type Auth struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type AuthSuccess struct {
}

func auth(msg *Auth) interface{} {
	return &AuthSuccess{}
}
