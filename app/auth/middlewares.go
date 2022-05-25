package auth

import (
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
)

func Check(c server.Context) error {
	var routeName = c.RouteName()
	var userRoles = []*entities.Role{ROLE_GUEST}
	var user = c.User()
	var authConfig = GetAuthConfig(routeName)

	// If there is no auth config for this route, then allow all
	if authConfig == nil {
		return c.Next()
	}

	if authConfig.Prepare != nil {
		if err := authConfig.Prepare(c); err != nil {
			if entities.IsNotFound(err) {
				return c.Status(http.StatusNotFound).SendString("Not found")
			}
			return err
		}
	}

	if user != nil && user.IsRoot() {
		return c.Next()
	}

	if user != nil {
		userRoles = user.Roles
	}

	if user.ID > 0 && !user.Active {
		c.Cookie(&server.Cookie{
			Name:    config.APP_TOKEN_KEY,
			Value:   "",
			Expires: time.Now().Add(time.Hour * 100 * 365 * 24),
		})

		return c.Redirect(utils.Url("/inactive"))
	}

	// Check all user roles for this action
	for _, role := range userRoles {
		permission := GetRolePermission(role.ID, routeName)

		if permission.Value == entities.PERM_ALL {
			return c.Next()
		}

		if permission.Value == entities.PERM_OWN && authConfig.OwnCheckFN != nil && authConfig.OwnCheckFN(c) {
			return c.Next()
		}
	}

	if user == nil || user.ID == 0 {
		return c.Redirect("/login?back=" + url.QueryEscape(c.OriginalURL()))
	}

	return c.Status(http.StatusForbidden).SendString("Insufficient permission")
}

func AssignUserInfo(c server.Context) error {
	c.Locals("user", GUEST_USER)
	tokenString := c.Cookies(config.APP_TOKEN_KEY)

	if tokenString == "" {
		return c.Next()
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&entities.UserJwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.APP_KEY), nil
		},
	)

	if err == nil {
		if claims, ok := token.Claims.(*entities.UserJwtClaims); ok && token.Valid {
			user := &claims.User
			user.Roles = GetRolesFromIDs(user.RoleIDs)
			c.Locals("user", user)
		}
	} else {
		c.Logger().Error(err)
	}

	return c.Next()
}
