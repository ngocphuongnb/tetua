package middlewares

import (
	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/server"
)

func All() []server.Handler {
	return []server.Handler{
		RequestID,
		auth.AssignUserInfo,
		Recover,
		RequestLog,
		Cookie,
		auth.Check,
	}
}
