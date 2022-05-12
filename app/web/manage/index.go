package manage

import (
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/views"
)

func Manage(c server.Context) (err error) {
	c.Meta().Title = "Manage"
	return c.Render(views.Manage())
}
