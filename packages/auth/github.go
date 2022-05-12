package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
	"golang.org/x/oauth2"
)

const GITHUB_ACCESS_TOKEN_URL = "https://github.com/login/oauth/access_token"
const GITHUB_USER_URL = "https://api.github.com/user"

type GithubAccessTokenResponse struct {
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

type GithubUserResponse struct {
	Login             string    `json:"login"`
	ID                int       `json:"id"`
	NodeID            string    `json:"node_id"`
	AvatarURL         string    `json:"avatar_url"`
	GravatarID        string    `json:"gravatar_id"`
	URL               string    `json:"url"`
	HTMLURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	ReposURL          string    `json:"repos_url"`
	EventsURL         string    `json:"events_url"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	Name              string    `json:"name"`
	Company           string    `json:"company"`
	Blog              string    `json:"blog"`
	Location          string    `json:"location"`
	Email             string    `json:"email"`
	Hireable          string    `json:"hireable"`
	Bio               string    `json:"bio"`
	TwitterUsername   string    `json:"twitter_username"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type GithubAuthProvider struct {
	config *oauth2.Config
}

func NewGithub(config *oauth2.Config) server.AuthProvider {
	if config.ClientID == "" || config.ClientSecret == "" {
		fmt.Println("Github client id or secret is not set")
		os.Exit(1)
	}
	return &GithubAuthProvider{
		config: config,
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

	githubId := strconv.Itoa(githubUser.ID)

	return repositories.User.CreateIfNotExistsByProvider(c.Context(), &entities.User{
		Provider:         "github",
		ProviderID:       utils.SanitizePlainText(githubId),
		Username:         utils.SanitizePlainText(githubUser.Login),
		Email:            utils.SanitizePlainText(githubUser.Email),
		ProviderAvatar:   utils.SanitizePlainText(githubUser.AvatarURL),
		DisplayName:      utils.SanitizePlainText(githubUser.Name),
		URL:              utils.SanitizePlainText(githubUser.Blog),
		ProviderUsername: utils.SanitizePlainText(githubUser.Login),
		RoleIDs:          []int{auth.ROLE_USER.ID},
		Active:           config.Setting("auto_approve_user") == "yes",
	})
}
