package repositories

import (
	"github.com/ngocphuongnb/tetua/app/entities"
)

type FileRepository interface {
	Repository[entities.File, entities.FileFilter]
}
