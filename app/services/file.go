package services

import (
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/fs"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/valyala/fasthttp"
)

func SaveFile(c server.Context, inputFileName string) (*entities.File, error) {
	featuredImageHeader, err := c.File(inputFileName)
	if err != nil {
		if err != fasthttp.ErrMissingFile {
			c.Logger().Error(err)
			return nil, err
		}
		return nil, nil
	}

	featuredImage, err := fs.Disk().PutMultipart(c.Context(), featuredImageHeader)
	if err != nil {
		return nil, err
	}

	return repositories.File.Create(c.Context(), &entities.File{
		Disk:   featuredImage.Disk,
		Path:   featuredImage.Path,
		Type:   featuredImage.Type,
		Size:   featuredImage.Size,
		UserID: c.User().ID,
	})
}
