package services_test

import (
	"fmt"
	"testing"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/fs"
	"github.com/ngocphuongnb/tetua/app/mock"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/services"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestFile(t *testing.T) {
	mock.CreateRepositories()
	fs.New("disk_mock", []fs.FSDisk{&mock.Disk{}})
	mockLogger := mock.CreateLogger(true)
	mockServer := mock.CreateServer()
	mockServer.Post("/test-upload-invalid-header", func(c server.Context) error {
		f, err := services.SaveFile(c, "featured_image")
		assert.Equal(t, (*entities.File)(nil), f)
		assert.Equal(t, fasthttp.ErrNoMultipartForm, err)
		assert.Equal(t, fasthttp.ErrNoMultipartForm.Error(), fmt.Sprintf("%v", mockLogger.Last().Params[0]))
		return c.SendString("ok")
	})

	mock.PostRequest(mockServer, "/test-upload-invalid-header")

	mockServer.Post("/test-upload-empty", func(c server.Context) error {
		f, err := services.SaveFile(c, "featured_image")
		assert.Equal(t, (*entities.File)(nil), f)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(mockLogger.Messages))
		return c.SendString("ok")
	})

	req := mock.CreateUploadRequest("POST", "/test-upload-empty", "some_field", "image.jpg")
	mockServer.Test(req)

	mockServer.Post("/test-upload-disk-error", func(c server.Context) error {
		f, err := services.SaveFile(c, "featured_image")
		assert.Equal(t, (*entities.File)(nil), f)
		assert.Equal(t, "PutMultipart error", err.Error())
		assert.Equal(t, 1, len(mockLogger.Messages))
		return c.SendString("ok")
	})

	req = mock.CreateUploadRequest("POST", "/test-upload-disk-error", "featured_image", "error.jpg")
	mockServer.Test(req)

	mockServer.Post("/test-upload-success", func(c server.Context) error {
		c.Locals("user", &entities.User{ID: 2})
		f, err := services.SaveFile(c, "featured_image")
		assert.NoError(t, err)
		assert.Equal(t, 1, f.ID)
		assert.Equal(t, "disk_mock", f.Disk)
		assert.Equal(t, "image.jpg", f.Path)
		assert.Equal(t, "image/jpeg", f.Type)
		assert.Equal(t, 100, f.Size)
		assert.Equal(t, 2, f.UserID)
		return c.SendString("ok")
	})

	req = mock.CreateUploadRequest("POST", "/test-upload-success", "featured_image", "image.jpg")
	mockServer.Test(req)
}
