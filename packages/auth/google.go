package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type GoogleUserResponse struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type GoogleAuthProvider struct {
	config *oauth2.Config
}

func NewGoogle(cfg map[string]string) server.AuthProvider {
	if cfg["client_id"] == "" || cfg["client_secret"] == "" {
		panic("Github client id or secret is not set")
	}

	return &GoogleAuthProvider{
		config: &oauth2.Config{
			ClientID:     cfg["client_id"],
			ClientSecret: cfg["client_secret"],
			Endpoint:     google.Endpoint,
			RedirectURL:  utils.Url("/auth/google/callback"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
		},
	}
}

func (g *GoogleAuthProvider) Name() string {
	return "google"
}

func (g *GoogleAuthProvider) GetGoogleUserFromAccessCode(code string) (*GoogleUserResponse, error) {
	token, err := g.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	userResponse := &GoogleUserResponse{}
	if err := json.Unmarshal(body, userResponse); err != nil {
		return nil, err
	}

	return userResponse, nil
}

func (g *GoogleAuthProvider) Login(c server.Context) error {
	url := g.config.AuthCodeURL(c.Cookies(config.COOKIE_UUID))

	return c.Redirect(url)
}

func (g *GoogleAuthProvider) Callback(c server.Context) (u *entities.User, err error) {
	if c.Query("state") != c.Cookies(config.COOKIE_UUID) {
		return nil, fmt.Errorf("invalid oauth google state")
	}

	if c.Query("code") == "" {
		return nil, fmt.Errorf("code is empty")
	}

	googleUser, err := g.GetGoogleUserFromAccessCode(c.Query("code"))

	if err != nil {
		return nil, err
	}

	return &entities.User{
		Provider:         "google",
		ProviderID:       utils.SanitizePlainText(googleUser.ID),
		Username:         utils.SanitizePlainText(strings.Split(googleUser.Email, "@gmail")[0]),
		Email:            utils.SanitizePlainText(googleUser.Email),
		ProviderAvatar:   utils.SanitizePlainText(googleUser.Picture),
		DisplayName:      utils.SanitizePlainText(googleUser.Name),
		URL:              utils.SanitizePlainText(""),
		ProviderUsername: utils.SanitizePlainText(googleUser.Email),
		RoleIDs:          []int{auth.ROLE_USER.ID},
		Active:           config.Setting("auto_approve_user") == "yes",
	}, nil
}
