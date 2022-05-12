package web_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/cache"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/fs"
	"github.com/ngocphuongnb/tetua/app/mock"
	mockrepository "github.com/ngocphuongnb/tetua/app/mock/repository"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/web"
	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/xml"
)

var m = minify.New()
var mockLogger *mock.MockLogger
var post1, post2 *entities.Post
var file1, file2, file3 *entities.File
var topic1 *entities.Topic

func init() {
	m.AddFunc("xml", xml.Minify)
	fs.New("disk_mock", []fs.FSDisk{&mock.Disk{}})
	cache.Roles = []*entities.Role{auth.ROLE_ADMIN, auth.ROLE_USER, auth.ROLE_GUEST}
	mockLogger = mock.CreateLogger(true)
	mock.CreateRepositories()
	config.Settings([]*config.SettingItem{{
		Name:  "app_base_url",
		Value: "http://localhost:8080",
	}})

	topic1, _ = repositories.Topic.Create(context.Background(), &entities.Topic{
		ID:   1,
		Name: "Test Topic",
		Slug: "test-topic",
	})

	post1, _ = repositories.Post.Create(context.Background(), &entities.Post{
		ID:       1,
		Name:     "test post 1",
		Slug:     "test-post-1",
		Draft:    false,
		UserID:   1,
		Topics:   []*entities.Topic{topic1},
		TopicIDs: []int{topic1.ID},
	})
	post2, _ = repositories.Post.Create(context.Background(), &entities.Post{
		ID:       2,
		Name:     "test post 2",
		Slug:     "test-post-2",
		Draft:    false,
		UserID:   2,
		Topics:   []*entities.Topic{topic1},
		TopicIDs: []int{topic1.ID},
	})

	file1, _ = repositories.File.Create(context.Background(), &entities.File{
		ID:     1,
		Disk:   "disk_mock",
		Path:   "/test/file1.jpg",
		Size:   100,
		Type:   "image/jpg",
		UserID: 1,
	})
	file2, _ = repositories.File.Create(context.Background(), &entities.File{
		ID:     2,
		Disk:   "disk_mock",
		Path:   "/test/file2.jpg",
		Size:   100,
		Type:   "image/jpg",
		UserID: 1,
	})
	file3, _ = repositories.File.Create(context.Background(), &entities.File{
		ID:     3,
		Disk:   "disk_mock",
		Path:   "/test/file3.jpg",
		Size:   100,
		Type:   "image/jpg",
		UserID: 2,
	})
}

func TestWeb(t *testing.T) {
	web.NewServer(web.Config{
		JwtSigningKey: config.APP_KEY,
		Theme:         config.APP_THEME,
	})
}

func TestFeedError(t *testing.T) {
	mockrepository.FakeRepoErrors["post_find"] = errors.New("Error finding posts")
	mockServer := mock.CreateServer()
	mockServer.Get("/feed", func(c server.Context) error {
		return web.Feed(c)
	})

	body, resp := mock.GetRequest(mockServer, "/feed")
	assert.Equal(t, "Error", body)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestFeed(t *testing.T) {
	mockrepository.FakeRepoErrors["post_find"] = nil
	mockServer := mock.CreateServer()
	mockServer.Get("/feed", func(c server.Context) error {
		return web.Feed(c)
	})

	body, _ := mock.GetRequest(mockServer, "/feed")
	body, err := m.String("xml", body)
	assert.Nil(t, err)

	expectFeed, _ := m.String("xml", fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<rss version="2.0"
		xmlns:content="http://purl.org/rss/1.0/modules/content/">
		<channel>
			<title>Tetua</title>
			<link>http://localhost:8080/</link>
			<description>Tetua</description>
			<item>
				<title>test post 1</title>
				<link>http://localhost:8080/test-post-1-1.html</link>
				<description/>
				<author>testuser1</author>
				<guid>1</guid>
				<pubDate>%s</pubDate>
			</item>
			<item>
				<title>test post 2</title>
				<link>http://localhost:8080/test-post-2-2.html</link>
				<description/>
				<author>testuser2</author>
				<guid>2</guid>
				<pubDate>%s</pubDate>
			</item>
		</channel>
	</rss>`, post1.CreatedAt.Format(time.RFC1123Z), post2.CreatedAt.Format(time.RFC1123Z)))

	assert.Equal(t, expectFeed, body)
}

func TestFileError(t *testing.T) {
	mockrepository.FakeRepoErrors["file_find"] = errors.New("Error finding files")
	mockServer := mock.CreateServer()
	mockServer.Get("/files", web.FileList)

	_, resp := mock.GetRequest(mockServer, "/files")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestFileList(t *testing.T) {
	mockrepository.FakeRepoErrors["file_find"] = nil

	fileLinks := []string{
		file1.Url(),
		file2.Url(),
	}

	mockServer := mock.CreateServer()
	mockServer.Get("/files", func(c server.Context) error {
		c.Locals("user", &entities.User{ID: 1})
		return web.FileList(c)
	})

	body, resp := mock.GetRequest(mockServer, "/files")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	assert.Nil(t, err)
	foundLinks := make([]string, 0)
	doc.Find(".files-list > div").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a").Attr("href")
		foundLinks = append(foundLinks, href)
	})

	assert.Equal(t, fileLinks, foundLinks)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFileDeleteError1(t *testing.T) {
	mockrepository.FakeRepoErrors["file_deleteByID"] = errors.New("Error deleting file")
	mockServer := mock.CreateServer()
	mockServer.Delete("/files/:id", func(c server.Context) error {
		file3, _ := repositories.File.ByID(context.Background(), 3)
		c.Locals("file", file3)
		return web.FileDelete(c)
	})

	body, resp := mock.Request(mockServer, "DELETE", "/files/3")
	assert.Equal(t, errors.New("Error deleting file"), mockLogger.Last().Params[0])
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"type":"error","message":"Error deleting file"}`, body)
}

func TestFileDeleteError2(t *testing.T) {
	mockrepository.FakeRepoErrors["file_deleteByID"] = nil
	mockServer := mock.CreateServer()
	mockServer.Delete("/files/:id", func(c server.Context) error {
		file3, _ := repositories.File.ByID(context.Background(), 3)
		file3.Path = "/delete/error"
		c.Locals("file", file3)
		return web.FileDelete(c)
	})

	body, resp := mock.Request(mockServer, "DELETE", "/files/3")
	assert.Equal(t, errors.New("Delete file error"), mockLogger.Last().Params[0])
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"type":"error","message":"Error deleting file"}`, body)
}

func TestFileDelete(t *testing.T) {
	mockrepository.FakeRepoErrors["file_deleteByID"] = nil
	repositories.File.Create(context.Background(), file3)
	mockServer := mock.CreateServer()
	mockServer.Delete("/files/:id", func(c server.Context) error {
		file3, _ := repositories.File.ByID(context.Background(), 3)
		file3.Path = "/test/file3.jpg"
		c.Locals("file", file3)
		return web.FileDelete(c)
	})

	body, resp := mock.Request(mockServer, "DELETE", "/files/3")

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `{"type":"success","message":"File deleted"}`, body)
}

func TestFileUpload(t *testing.T) {
	mockServer := mock.CreateServer()
	mockServer.Post("/files", func(c server.Context) error {
		return web.Upload(c)
	})

	req := mock.CreateUploadRequest("POST", "/files", "some_field", "file.jpg")
	body, _ := mock.SendRequest(mockServer, req)
	assert.Equal(t, errors.New("there is no uploaded file associated with the given key"), mockLogger.Last().Params[0])
	assert.Equal(t, `{"error":"Error saving file"}`, body)

	req = mock.CreateUploadRequest("POST", "/files", "file", "error.jpg")
	body, _ = mock.SendRequest(mockServer, req)
	assert.Equal(t, errors.New("PutMultipart error"), mockLogger.Last().Params[0])
	assert.Equal(t, `{"error":"Error saving file"}`, body)

	mockrepository.FakeRepoErrors["file_create"] = errors.New("Error creating file")
	req = mock.CreateUploadRequest("POST", "/files", "file", "image.jpg")
	body, _ = mock.SendRequest(mockServer, req)
	assert.Equal(t, errors.New("Error creating file"), mockLogger.Last().Params[0])
	assert.Equal(t, `{"error":"Error saving file"}`, body)

	mockrepository.FakeRepoErrors["file_create"] = nil
	req = mock.CreateUploadRequest("POST", "/files", "file", "image.jpg")
	body, _ = mock.SendRequest(mockServer, req)
	assert.Equal(t, `{"size":100,"type":"image/jpeg","url":"/disk_mock/image.jpg"}`, body)
}

func TestIndex(t *testing.T) {
	mockServer := mock.CreateServer()
	mockServer.Get("/", func(c server.Context) error {
		return web.Index(c)
	})

	mockrepository.FakeRepoErrors["post_paginate"] = errors.New("Error paginating posts")
	body, resp := mock.GetRequest(mockServer, "/")
	assert.Equal(t, http.StatusBadGateway, resp.StatusCode)
	assert.Equal(t, errors.New("Error paginating posts"), mockLogger.Last().Params[0])
	assert.Equal(t, true, strings.Contains(body, `<h1>Something went wrong</h1>`))

	mockrepository.FakeRepoErrors["post_paginate"] = nil
	mockrepository.FakeRepoErrors["post_find"] = errors.New("Error finding post")
	body, resp = mock.GetRequest(mockServer, "/")
	assert.Equal(t, http.StatusBadGateway, resp.StatusCode)
	assert.Equal(t, errors.New("Error finding post"), mockLogger.Last().Params[0])
	assert.Equal(t, true, strings.Contains(body, `<h1>Something went wrong</h1>`))

	mockrepository.FakeRepoErrors["post_paginate"] = nil
	mockrepository.FakeRepoErrors["post_find"] = nil
	body, resp = mock.GetRequest(mockServer, "/")

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	assert.Nil(t, err)
	postLinks := []string{post1.Url(), post2.Url()}
	mainLinks := make([]string, 0)
	sideLinks := make([]string, 0)
	doc.Find("main article").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a.overlay").Attr("href")
		mainLinks = append(mainLinks, href)
	})

	doc.Find(".right article").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a").Attr("href")
		sideLinks = append(sideLinks, href)
	})

	assert.Equal(t, postLinks, mainLinks)
	assert.Equal(t, postLinks, sideLinks)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSearch(t *testing.T) {
	mockServer := mock.CreateServer()
	mockServer.Get("/search", func(c server.Context) error {
		return web.Search(c)
	})

	mockrepository.FakeRepoErrors["post_paginate"] = errors.New("Error paginating posts")
	body, resp := mock.GetRequest(mockServer, "/search?q=post")
	assert.Equal(t, http.StatusBadGateway, resp.StatusCode)
	assert.Equal(t, errors.New("Error paginating posts"), mockLogger.Last().Params[0])
	assert.Equal(t, true, strings.Contains(body, `<h1>Something went wrong</h1>`))

	mockrepository.FakeRepoErrors["post_paginate"] = nil
	mockrepository.FakeRepoErrors["post_find"] = nil
	body, resp = mock.GetRequest(mockServer, "/search?q=post")

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	assert.Nil(t, err)
	postLinks := []string{post1.Url(), post2.Url()}
	mainLinks := make([]string, 0)
	doc.Find("main article").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a.overlay").Attr("href")
		mainLinks = append(mainLinks, href)
	})

	assert.Equal(t, postLinks, mainLinks)
	assert.Equal(t, true, strings.Contains(body, `<title>post - Search result for post - Tetua</title>`))
}

func TestTopicView(t *testing.T) {
	mockServer := mock.CreateServer()
	mockServer.Get("/:slug", func(c server.Context) error {
		return web.TopicView(c)
	})

	body, resp := mock.GetRequest(mockServer, "/test-topic")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, true, strings.Contains(body, `<title>Topic not found`))
	assert.Equal(t, true, strings.Contains(body, `<h1>Topic not found</h1>`))

	post1.TopicIDs = []int{topic1.ID}
	post2.TopicIDs = []int{topic1.ID}
	repositories.Post.Update(context.Background(), post1)
	repositories.Post.Update(context.Background(), post2)
	cache.CacheTopics()

	mockrepository.FakeRepoErrors["post_paginate"] = errors.New("Error paginating posts")
	body, resp = mock.GetRequest(mockServer, "/test-topic")
	assert.Equal(t, http.StatusBadGateway, resp.StatusCode)
	assert.Equal(t, errors.New("Error paginating posts"), mockLogger.Last().Params[0])
	assert.Equal(t, true, strings.Contains(body, `<h1>Something went wrong</h1>`))

	mockrepository.FakeRepoErrors["post_paginate"] = nil
	mockrepository.FakeRepoErrors["post_find"] = errors.New("Error finding post")
	body, resp = mock.GetRequest(mockServer, "/test-topic")
	assert.Equal(t, http.StatusBadGateway, resp.StatusCode)
	assert.Equal(t, errors.New("Error finding post"), mockLogger.Last().Params[0])
	assert.Equal(t, true, strings.Contains(body, `<h1>Something went wrong</h1>`))

	mockrepository.FakeRepoErrors["post_find"] = nil
	body, resp = mock.GetRequest(mockServer, "/test-topic")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, true, strings.Contains(body, `<title>Test Topic`))

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	assert.Nil(t, err)
	postLinks := []string{post1.Url(), post2.Url()}
	mainLinks := make([]string, 0)
	sideLinks := make([]string, 0)
	doc.Find("main article").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a.overlay").Attr("href")
		mainLinks = append(mainLinks, href)
	})

	doc.Find(".right article").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a").Attr("href")
		sideLinks = append(sideLinks, href)
	})

	assert.Equal(t, postLinks, mainLinks)
	assert.Equal(t, postLinks, sideLinks)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestTopicFeedError(t *testing.T) {
	mockrepository.FakeRepoErrors["post_find"] = errors.New("Error finding posts")
	mockServer := mock.CreateServer()
	mockServer.Get("/:slug/feed", func(c server.Context) error {
		return web.TopicFeed(c)
	})

	topic1.Slug = "test-topic-updated"
	repositories.Topic.Update(context.Background(), topic1)
	cache.CacheTopics()
	body, resp := mock.GetRequest(mockServer, "/test-topic/feed")
	assert.Equal(t, "Topic not found", body)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	topic1.Slug = "test-topic"
	mockrepository.FakeRepoErrors["post_find"] = nil
	repositories.Topic.Update(context.Background(), topic1)
	cache.CacheTopics()
	body, resp = mock.GetRequest(mockServer, "/test-topic/feed")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := m.String("xml", body)
	assert.Nil(t, err)

	expectFeed, _ := m.String("xml", fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<rss version="2.0"
		xmlns:content="http://purl.org/rss/1.0/modules/content/">
		<channel>
			<title>%s</title>
			<link>%s</link>
			<description></description>
			<item>
				<title>%s</title>
				<link>%s</link>
				<description></description>
				<author>%s</author>
				<guid>%d</guid>
				<pubDate>%s</pubDate>
			</item>
			<item>
				<title>%s</title>
				<link>%s</link>
				<description></description>
				<author>%s</author>
				<guid>%d</guid>
				<pubDate>%s</pubDate>
			</item>
		</channel>
	</rss>`,
		topic1.Name,
		topic1.Url(),

		post1.Name,
		post1.Url(),
		post1.User.Username,
		post1.ID,
		post1.CreatedAt.Format(time.RFC1123Z),

		post2.Name,
		post2.Url(),
		post2.User.Username,
		post2.ID,
		post2.CreatedAt.Format(time.RFC1123Z),
	))
	assert.Equal(t, expectFeed, body)
}
