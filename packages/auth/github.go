package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const GITHUB_ACCESS_TOKEN_URL = "https://github.com/login/oauth/access_token"
const GITHUB_USER_URL = "https://api.github.com/user"

type GithubAccessTokenResponse struct {
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

type GithubUserResponse struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Blog      string `json:"blog"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
}

type GithubAuthProvider struct {
	config *oauth2.Config
}

func NewGithub(cfg map[string]string) server.AuthProvider {
	if cfg["client_id"] == "" || cfg["client_secret"] == "" {
		panic("Github client id or secret is not set")
	}
	return &GithubAuthProvider{
		config: &oauth2.Config{
			ClientID:     cfg["client_id"],
			ClientSecret: cfg["client_secret"],
			RedirectURL:  utils.Url("/auth/github/callback"),
			Endpoint:     github.Endpoint,
		},
	}
}

func (g *GithubAuthProvider) Name() string {
	return "github"
}

func (g *GithubAuthProvider) GetGithubAccessToken(code string) (string, error) {
	requestBody := map[string]string{
		"code":          code,
		"client_id":     g.config.ClientID,
		"client_secret": g.config.ClientSecret,
	}
	requestJSON, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(
		"POST",
		GITHUB_ACCESS_TOKEN_URL,
		bytes.NewBuffer(requestJSON),
	)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var accessTokenResponse GithubAccessTokenResponse
	if err := json.Unmarshal(body, &accessTokenResponse); err != nil {
		return "", err
	}

	return accessTokenResponse.AccessToken, nil
}

func (g *GithubAuthProvider) GetGithubUser(accessToken string) (*GithubUserResponse, error) {
	req, err := http.NewRequest(
		"GET",
		GITHUB_USER_URL,
		nil,
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	userResponse := &GithubUserResponse{}
	if err := json.Unmarshal(body, userResponse); err != nil {
		return nil, err
	}

	return userResponse, nil
}

func (g *GithubAuthProvider) GetGithubUserFromAccessCode(code string) (*GithubUserResponse, error) {
	accessToken, err := g.GetGithubAccessToken(code)
	if err != nil {
		return nil, err
	}

	return g.GetGithubUser(accessToken)
}

func (g *GithubAuthProvider) Login(c server.Context) error {
	url := g.config.AuthCodeURL(
		c.Cookies(config.COOKIE_UUID),
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("scope", "user:email"),
	)

	return c.Redirect(url)
}

func (g *GithubAuthProvider) Callback(c server.Context) (u *entities.User, err error) {
	if c.Query("code") == "" {
		return nil, fmt.Errorf("code is empty")
	}

	githubUser, err := g.GetGithubUserFromAccessCode(c.Query("code"))

	if err != nil {
		return nil, err
	}

	return &entities.User{
		Provider:         "github",
		ProviderID:       utils.SanitizePlainText(strconv.Itoa(githubUser.ID)),
		Username:         utils.SanitizePlainText("github_" + strconv.Itoa(githubUser.ID) + "_" + githubUser.Login),
		Email:            utils.SanitizePlainText(githubUser.Email),
		ProviderAvatar:   utils.SanitizePlainText(githubUser.AvatarURL),
		DisplayName:      utils.SanitizePlainText(githubUser.Name),
		URL:              utils.SanitizePlainText(githubUser.Blog),
		ProviderUsername: utils.SanitizePlainText(githubUser.Login),
		RoleIDs:          []int{auth.ROLE_USER.ID},
		Active:           config.Setting("auto_approve_user") == "yes",
	}, nil
}
