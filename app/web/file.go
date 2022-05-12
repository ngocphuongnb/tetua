package web

import (
	"net/http"
	"sync"

	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/fs"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/views"
)

func FileList(c server.Context) error {
	paginate, err := repositories.File.Paginate(c.Context(), &entities.FileFilter{
		UserIDs: []int{c.User().ID},
		Filter: &entities.Filter{
			BaseUrl:         config.Url("/files"),
			Page:            c.QueryInt("page"),
			Limit:           12,
			IgnoreUrlParams: []string{"user"},
		},
	})

	if err != nil {
		c.Logger().Error(err)
		return c.Status(http.StatusInternalServerError).Render(views.Error("Something went wrong"))
	}

	return c.Render(views.FileList(paginate))
}

func FileDelete(c server.Context) (err error) {
	var wg sync.WaitGroup
	var err1, err2 error
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err1 = repositories.File.DeleteByID(c.Context(), c.ParamInt("id"))
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		if file, ok := c.Locals("file").(*entities.File); ok && file != nil {
			err2 = file.Delete(c.Context())
		}
	}(&wg)

	wg.Wait()

	if err := utils.FirstError(err1, err2); err != nil {
		c.Logger().Error(err)
		return c.Status(http.StatusBadRequest).Json(&entities.Message{
			Type:    "error",
			Message: "Error deleting file",
		})
	}

	return c.Status(http.StatusOK).Json(&entities.Message{
		Type:    "success",
		Message: "File deleted",
	})
}

func Upload(c server.Context) error {
	if uploadFile, err := c.File("file"); err == nil {
		if uploadedFile, err := fs.Disk().PutMultipart(c.Context(), uploadFile); err != nil {
			c.Logger().Error(err)
		} else {
			f, err := repositories.File.Create(c.Context(), &entities.File{
				Disk:   uploadedFile.Disk,
				Path:   uploadedFile.Path,
				Type:   uploadedFile.Type,
				Size:   uploadedFile.Size,
				UserID: c.User().ID,
			})

			if err != nil {
				c.Logger().Error(err)
				return c.Status(http.StatusInternalServerError).Json(entities.Map{
					"error": "Error saving file",
				})
			}

			return c.Json(entities.Map{
				"size": f.Size,
				"type": f.Type,
				"url":  f.Url(),
			})
		}
	} else {
		c.Logger().Error(err)
	}

	return c.Status(http.StatusInternalServerError).Json(entities.Map{
		"error": "Error saving file",
	})
}
