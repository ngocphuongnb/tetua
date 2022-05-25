package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
)

const oauthTwitterUrlAPI = "https://api.twitter.com/1.1/account/verify_credentials.json"

type TwitterUserResponse struct {
	ID                   int    `json:"id"`
	IDStr                string `json:"id_str"`
	Name                 string `json:"name"`
	ScreenName           string `json:"screen_name"`
	Email                string `json:"email"`
	ProfileImageUrlHttps string `json:"profile_image_url_https"`
}

type TwitterAuthProvider struct {
	config *oauth1.Config
}

func NewTwitter(cfg map[string]string) server.AuthProvider {
	if cfg["consumer_key"] == "" || cfg["consumer_secret"] == "" {
		panic("Github client id or secret is not set")
	}

	return &TwitterAuthProvider{
		config: &oauth1.Config{
			ConsumerKey:    cfg["consumer_key"],
			ConsumerSecret: cfg["consumer_secret"],
			CallbackURL:    utils.Url("/auth/twitter/callback"),
			Endpoint:       twitter.AuthorizeEndpoint,
		},
	}
}

func (g *TwitterAuthProvider) Name() string {
	return "twitter"
}

func (g *TwitterAuthProvider) Login(c server.Context) error {
	requestToken, _, err := g.config.RequestToken()
	url, err := g.config.AuthorizationURL(requestToken)

	if err != nil {
		return err
	}

	return c.Redirect(url.String())
}

func (g *TwitterAuthProvider) Callback(c server.Context) (u *entities.User, err error) {
	requestToken := c.Query("oauth_token")
	verifier := c.Query("oauth_verifier")

	if requestToken == "" || verifier == "" {
		return nil, errors.New("oauth1: Request missing oauth_token or oauth_verifier")
	}

	accessToken, accessSecret, err := g.config.AccessToken(requestToken, "", verifier)

	if err != nil {
		return nil, err
	}

	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := g.config.Client(oauth1.NoContext, token)
	resp, err := httpClient.Get(oauthTwitterUrlAPI)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	userResponse := &TwitterUserResponse{}
	if err := json.Unmarshal(body, userResponse); err != nil {
		return nil, err
	}

	return &entities.User{
		Provider:         "twitter",
		ProviderID:       utils.SanitizePlainText(userResponse.IDStr),
		Username:         utils.SanitizePlainText("twitter_" + userResponse.IDStr + "_" + userResponse.ScreenName),
		Email:            utils.SanitizePlainText(userResponse.Email),
		ProviderAvatar:   utils.SanitizePlainText(userResponse.ProfileImageUrlHttps),
		DisplayName:      utils.SanitizePlainText(userResponse.Name),
		URL:              utils.SanitizePlainText(""),
		ProviderUsername: utils.SanitizePlainText(userResponse.ScreenName),
		RoleIDs:          []int{auth.ROLE_USER.ID},
		Active:           config.Setting("auto_approve_user") == "yes",
	}, nil
}
