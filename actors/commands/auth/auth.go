package auth

import "github.com/AsynkronIT/gam/actor"

type Auth struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type AuthSuccess struct {
}

type AuthFail struct {
}

type authActor struct {
}

func NewActor() actor.Actor {
	return &authActor{}
}

func (state *authActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Auth:
		if err := fetch(msg); err != nil {
		}
		ctx.Respond(&AuthSuccess{})
	}
}

func fetch(auth *Auth) error {
	return nil
}
