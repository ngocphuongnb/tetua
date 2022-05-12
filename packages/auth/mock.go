package auth

import (
	"fmt"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/server"
	"golang.org/x/oauth2"
)

type MockAuthProvider struct {
	config *oauth2.Config
}

var MockRootUser = &entities.User{
	ID:               1,
	Username:         "mockrootuser",
	Email:            "mockrootuser@mock-service.local",
	Provider:         "mock",
	ProviderID:       "1",
	ProviderUsername: "mockrootuser",
	RoleIDs:          []int{1},
	Roles:            []*entities.Role{{ID: 1}},
}

var MockNormalUser2 = &entities.User{
	ID:               2,
	Username:         "mocknormaluser1",
	Email:            "mocknormaluser1@mock-service.local",
	Provider:         "mock",
	ProviderID:       "2",
	ProviderUsername: "mocknormaluser1",
	RoleIDs:          []int{2},
	Roles:            []*entities.Role{{ID: 2}},
	Active:           true,
}

var MockNormalUser3 = &entities.User{
	ID:               3,
	Username:         "mocknormaluser2",
	Email:            "mocknormaluser2@mock-service.local",
	Provider:         "mock",
	ProviderID:       "3",
	ProviderUsername: "mockuser",
	RoleIDs:          []int{2},
	Roles:            []*entities.Role{{ID: 2}},
	Active:           true,
}

func NewMock(config *oauth2.Config) server.AuthProvider {
	return &MockAuthProvider{
		config: config,
	}
}

func (g *MockAuthProvider) Name() string {
	return "mock"
}

func (g *MockAuthProvider) Login(c server.Context) error {
	return c.Redirect("https://auth.mock-service.local/auth")
}

func (g *MockAuthProvider) Callback(c server.Context) (u *entities.User, err error) {
	if c.Query("code") == "" {
		return nil, fmt.Errorf("code is empty")
	}

	return MockNormalUser2, nil
}
