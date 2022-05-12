package webuser_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ngocphuongnb/tetua/app/mock"
	webuser "github.com/ngocphuongnb/tetua/app/web/user"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	mockServer := mock.CreateServer()
	mockServer.Get("/login", webuser.Login)
	body, resp := mock.GetRequest(mockServer, "/login")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, true, strings.Contains(body, `<h1 class="text-center">Login</h1>`))
}
