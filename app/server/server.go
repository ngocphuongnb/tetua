package server

import (
	"bufio"
	"context"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/logger"
)

type AuthProvider interface {
	Name() string
	Login(Context) error
	Callback(Context) (*entities.User, error)
}

type AuthConfig struct {
	Action       string
	Value        entities.PermType
	DefaultValue entities.PermType
	Prepare      func(c Context) error
	OwnCheckFN   func(c Context) bool
}

type Cookie struct {
	Name        string    `json:"name"`
	Value       string    `json:"value"`
	Path        string    `json:"path"`
	Domain      string    `json:"domain"`
	MaxAge      int       `json:"max_age"`
	Expires     time.Time `json:"expires"`
	Secure      bool      `json:"secure"`
	HTTPOnly    bool      `json:"http_only"`
	SameSite    string    `json:"same_site"`
	SessionOnly bool      `json:"session_only"`
}

type StaticConfig struct {
	Compress      bool          `json:"compress"`
	ByteRange     bool          `json:"byte_range"`
	Browse        bool          `json:"browse"`
	Download      bool          `json:"download"`
	Index         string        `json:"index"`
	CacheDuration time.Duration `json:"cache_duration"` // Default value 10 * time.Second.
	MaxAge        int           `json:"max_age"`        // Default value 0
}

type Context interface {
	RequestID() string
	Hostname() string
	BaseUrl() string
	RouteName() string
	User() *entities.User
	Post(...*entities.Post) *entities.Post
	Meta(meta ...*entities.Meta) *entities.Meta
	Messages(...*entities.Messages) *entities.Messages
	Logger() logger.Logger
	Json(interface{}) error
	Cookie(*Cookie)
	Cookies(string, ...string) string
	Next() error
	Status(int) Context
	Locals(string, ...interface{}) (val interface{})
	Response() Response
	Method() string
	OriginalURL() string
	Path() string
	IP() string
	Header(string, ...string) string
	Param(string) string
	ParamInt(string, ...int) int
	QueryInt(string, ...int) int
	Query(string, ...string) string
	SendString(string) error
	Send([]byte) error
	Redirect(string) error
	RedirectToRoute(name string, params ...map[string]interface{}) error
	BodyParser(interface{}) error
	Render(func(meta *entities.Meta, wr *bufio.Writer)) error
	Context() context.Context
	File(name string) (*multipart.FileHeader, error)
	WithError(msg string, err error)
}

type Handler func(c Context) error

type Server interface {
	Test(*http.Request, ...int) (*http.Response, error)
	Listen(string)
	Static(string, string, ...StaticConfig)
	Register(func(s Server)) Server
	Use(...Handler)
	UsePrefix(string, ...Handler)
	Group(string, ...Handler) Group
	Get(string, Handler, ...*AuthConfig)
	Post(string, Handler, ...*AuthConfig)
	Put(string, Handler, ...*AuthConfig)
	Delete(string, Handler, ...*AuthConfig)
	Patch(string, Handler, ...*AuthConfig)
	Head(string, Handler, ...*AuthConfig)
	Options(string, Handler, ...*AuthConfig)
}

type Response interface {
	StatusCode() int
	Header(string, ...string) string
}

type Group interface {
	Use(Handler)
	Group(string, ...Handler) Group
	Get(string, Handler, ...*AuthConfig)
	Post(string, Handler, ...*AuthConfig)
	Put(string, Handler, ...*AuthConfig)
	Delete(string, Handler, ...*AuthConfig)
	Patch(string, Handler, ...*AuthConfig)
	Head(string, Handler, ...*AuthConfig)
	Options(string, Handler, ...*AuthConfig)
}

type Route interface {
	Path() string
	Method() string
	Name() string
}
