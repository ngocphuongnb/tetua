package webpost

import (
	"net/http"

	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/views"
)

func ViewPage(c server.Context) error {
	page, err := repositories.Page.PublishedPageBySlug(c.Context(), c.Param("slug"))

	if err != nil {
		return c.Status(http.StatusNotFound).Render(views.Error("Page not found"))
	}

	c.Meta().Title = page.Name
	c.Meta().Description = page.Name

	return c.Render(views.PageView(page))
}
