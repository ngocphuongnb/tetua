package mock

import (
	"fmt"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/server"
	"golang.org/x/oauth2"
)

type AuthProvider struct {
	config *oauth2.Config
}

var RootUser = &entities.User{
	ID:               1,
	Username:         "mockrootuser",
	Email:            "mockrootuser@mock-service.local",
	Provider:         "mock",
	ProviderID:       "1",
	ProviderUsername: "mockrootuser",
	RoleIDs:          []int{1},
	Roles:            []*entities.Role{{ID: 1}},
}

var NormalUser2 = &entities.User{
	ID:               2,
	Username:         "mocknormaluser2",
	Email:            "mocknormaluser2@mock-service.local",
	Provider:         "mock",
	ProviderID:       "2",
	ProviderUsername: "mocknormaluser2",
	RoleIDs:          []int{2},
	Roles:            []*entities.Role{{ID: 2}},
	Active:           true,
}

var NormalUser3 = &entities.User{
	ID:               3,
	Username:         "mocknormaluser3",
	Email:            "mocknormaluser3@mock-service.local",
	Provider:         "mock",
	ProviderID:       "3",
	ProviderUsername: "mockuser",
	RoleIDs:          []int{2},
	Roles:            []*entities.Role{{ID: 2}},
	Active:           true,
}

func NewAuth(cfg map[string]string) server.AuthProvider {
	return &AuthProvider{
		config: &oauth2.Config{},
	}
}

func (g *AuthProvider) Name() string {
	return "mock"
}

func (g *AuthProvider) Login(c server.Context) error {
	return c.Redirect("https://auth.mock-service.local/auth")
}

func (g *AuthProvider) Callback(c server.Context) (u *entities.User, err error) {
	if c.Query("code") == "" {
		return nil, fmt.Errorf("code is empty")
	}

	return NormalUser2, nil
}
