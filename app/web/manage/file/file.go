package managefile

import (
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/views"
)

func Index(c server.Context) error {
	paginate, err := repositories.File.Paginate(c.Context(), &entities.FileFilter{
		Filter: &entities.Filter{
			BaseUrl: config.Url("/manage/files"),
			Page:    c.QueryInt("page"),
			Limit:   24,
		},
	})

	if err != nil {
		c.WithError("Something went wrong", err)
	}

	return c.Render(views.ManageFileIndex(paginate))
}
