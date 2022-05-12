package middlewares

import (
	"time"

	"github.com/google/uuid"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/server"
)

func Cookie(c server.Context) error {
	if c.Cookies(config.COOKIE_UUID) == "" {
		exp := time.Now().Add(time.Hour * 100 * 365 * 24)
		c.Cookie(&server.Cookie{
			Name:     config.COOKIE_UUID,
			Value:    uuid.NewString(),
			Expires:  exp,
			HTTPOnly: false,
			SameSite: "lax",
			Secure:   true,
		})
	}
	return c.Next()
}
