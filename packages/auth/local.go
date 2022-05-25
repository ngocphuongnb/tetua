package auth

import (
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
)

type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthAuthProvider struct {
}

func NewLocal(cfg map[string]string) server.AuthProvider {
	return &AuthAuthProvider{}
}

func (g *AuthAuthProvider) Name() string {
	return "local"
}

func (g *AuthAuthProvider) Login(c server.Context) error {
	return c.Redirect(utils.Url("/login"))
}
func (g *AuthAuthProvider) Callback(c server.Context) (*entities.User, error) {
	return nil, nil
}
