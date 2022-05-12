package auth

import (
	"net/http"
	"time"

	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/server"
)

var ActionConfigs = []*server.AuthConfig{}
var ROLE_ADMIN = &entities.Role{
	ID:   1,
	Name: "Admin",
	Root: true,
}

var ROLE_USER = &entities.Role{
	ID:   2,
	Name: "User",
	Root: false,
}

var ROLE_GUEST = &entities.Role{
	ID:   3,
	Name: "Guest",
	Root: false,
}

var GUEST_USER = &entities.User{
	ID:       0,
	Username: "Guest",
	Roles:    []*entities.Role{ROLE_GUEST},
}

func Config(cfg *server.AuthConfig) *server.AuthConfig {
	for _, ActionConfig := range ActionConfigs {
		if ActionConfig.Action == cfg.Action {
			panic("Duplicate action config: " + ActionConfig.Action)
		}
	}

	ActionConfigs = append(ActionConfigs, cfg)

	return cfg
}

func GetAuthConfig(action string) *server.AuthConfig {
	for _, config := range ActionConfigs {
		if config.Action == action {
			return config
		}
	}

	return nil
}

func SetLoginInfo(c server.Context, user *entities.User) error {
	exp := time.Now().Add(time.Hour * 100 * 365 * 24)
	jwtHeader, _ := c.Locals("jwt_header").(map[string]interface{})
	jwtToken, err := user.JwtClaim(exp, jwtHeader)

	if err != nil {
		return err
	}

	c.Cookie(&server.Cookie{
		Name:     config.APP_TOKEN_KEY,
		Value:    jwtToken,
		Expires:  exp,
		HTTPOnly: false,
		SameSite: "lax",
		Secure:   true,
	})

	return nil
}

func Routes(s server.Server) {
	authRoute := s.Group("/auth/:provider", func(c server.Context) error {
		provider := c.Param("provider")

		if GetProvider(provider) == nil {
			c.Status(http.StatusNotFound)
			return c.SendString("Invalid provider")
		}

		return c.Next()
	})

	authRoute.Get("", func(c server.Context) error {
		provider := GetProvider(c.Param("provider"))
		return provider.Login(c)
	})

	authRoute.Get("/callback", func(c server.Context) error {
		provider := GetProvider(c.Param("provider"))
		user, err := provider.Callback(c)

		if err != nil {
			c.Logger().Error(err)
			return c.Status(http.StatusBadGateway).SendString("Something went wrong")
		}

		if err = SetLoginInfo(c, user); err != nil {
			c.Logger().Error("Error setting login info", err)
			return c.Status(http.StatusBadGateway).SendString("Something went wrong")
		}

		return c.Redirect("/")
	})
}
