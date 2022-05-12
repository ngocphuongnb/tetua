package entities_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/fs"
	"github.com/ngocphuongnb/tetua/app/mock"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/stretchr/testify/assert"
)

func TestEntities(t *testing.T) {
	notfoundErr := &entities.NotFoundError{Message: "Not found"}
	assert.Equal(t, "Not found", notfoundErr.Error())
	assert.Equal(t, true, entities.IsNotFound(notfoundErr))
	assert.Equal(t, false, entities.IsNotFound(nil))

	messages := &entities.Messages{}
	assert.Equal(t, false, messages.HasError())
	assert.Equal(t, 0, messages.Length())
	messages.Append(&entities.Message{Type: "info", Message: "Info message"})
	messages.AppendError("Error message")

	assert.Equal(t, 2, len(*messages))
	assert.Equal(t, 2, messages.Length())
	assert.Equal(t, true, messages.HasError())

	assert.Equal(t, []*entities.Message{{
		Type:    "info",
		Message: "Info message",
	}, {
		Type:    "error",
		Message: "Error message",
	}}, messages.Get())

	sort := []*entities.Sort{{
		Field: "id",
		Order: "asc",
	}}
	filter := &entities.Filter{
		BaseUrl:         "/test",
		Search:          "search",
		Limit:           10,
		Page:            2,
		IgnoreUrlParams: []string{"some"},
		ExcludeIDs:      []int{1, 2},
		Sorts:           sort,
	}

	assert.Equal(t, "search", filter.GetSearch())
	assert.Equal(t, 10, filter.GetLimit())
	assert.Equal(t, 2, filter.GetPage())
	assert.Equal(t, sort, filter.GetSorts())
	assert.Equal(t, []string{"some"}, filter.GetIgnoreUrlParams())
	assert.Equal(t, []int{1, 2}, filter.GetExcludeIDs())
	assert.Equal(t, "/test?q=search", filter.Base())
	filter.BaseUrl = ""
	assert.Equal(t, "/?q=search", filter.Base())

	meta1 := &entities.Meta{}
	assert.Equal(t, config.Setting("app_name"), meta1.GetTitle())

	meta2 := &entities.Meta{Title: "Test"}
	assert.Equal(t, fmt.Sprintf("Test - %s", config.Setting("app_name")), meta2.GetTitle())

	paginate := &entities.Paginate[entities.Post]{
		Total:       100,
		PageSize:    10,
		PageCurrent: 2,
		Data:        utils.Repeat(&entities.Post{}, 10),
	}

	links := paginate.Links()

	for i, link := range links {
		page := i + 1
		expect := &entities.PaginateLink{Link: fmt.Sprintf("?page=%d", page), Label: fmt.Sprint(page), Class: ""}

		if page == 1 {
			expect = &entities.PaginateLink{Link: fmt.Sprintf("?page=%d", page), Label: fmt.Sprint(page), Class: "first"}
		}
		if page == 2 {
			expect = &entities.PaginateLink{Link: fmt.Sprintf("?page=%d", page), Label: fmt.Sprint(page), Class: "active"}
		}
		if page == len(links) {
			expect = &entities.PaginateLink{Link: fmt.Sprintf("?page=%d", page), Label: fmt.Sprint(page), Class: "last"}
		}

		assert.Equal(t, expect, link)
	}

	paginate.BaseUrl = "/test?tab=home"
	links = paginate.Links()

	for i, link := range links {
		page := i + 1
		expect := &entities.PaginateLink{Link: fmt.Sprintf("/test?tab=home&page=%d", page), Label: fmt.Sprint(page), Class: ""}

		if page == 1 {
			expect = &entities.PaginateLink{Link: fmt.Sprintf("/test?tab=home&page=%d", page), Label: fmt.Sprint(page), Class: "first"}
		}
		if page == 2 {
			expect = &entities.PaginateLink{Link: fmt.Sprintf("/test?tab=home&page=%d", page), Label: fmt.Sprint(page), Class: "active"}
		}
		if page == len(links) {
			expect = &entities.PaginateLink{Link: fmt.Sprintf("/test?tab=home&page=%d", page), Label: fmt.Sprint(page), Class: "last"}
		}

		assert.Equal(t, expect, link)
	}
}

func TestComment(t *testing.T) {
	commentFilter := &entities.CommentFilter{
		Filter: &entities.Filter{
			BaseUrl:         "/comment",
			Search:          "test",
			Page:            2,
			Limit:           10,
			Sorts:           []*entities.Sort{{"created_at", "desc"}},
			IgnoreUrlParams: []string{},
			ExcludeIDs:      []int{1, 2},
		},
		PostIDs:   []int{1, 2},
		UserIDs:   []int{1, 2},
		ParentIDs: []int{1, 2},
	}

	assert.Equal(t, "/comment?parent=1&post=1&q=test&user=1", commentFilter.Base())
	commentFilter.IgnoreUrlParams = []string{"parent", "user"}
	assert.Equal(t, "/comment?post=1&q=test", commentFilter.Base())
}

func TestFile(t *testing.T) {
	fileFilter := &entities.FileFilter{
		Filter: &entities.Filter{
			BaseUrl:         "/file",
			Search:          "test",
			Page:            2,
			Limit:           10,
			Sorts:           []*entities.Sort{{"created_at", "desc"}},
			IgnoreUrlParams: []string{},
			ExcludeIDs:      []int{1, 2},
		},
		UserIDs: []int{1, 2},
	}

	assert.Equal(t, "/file?q=test&user=1", fileFilter.Base())

	fs.New("disk_mock", []fs.FSDisk{&mock.Disk{}})
	file := &entities.File{
		ID:   1,
		Disk: "disk_mock",
		Path: "test/file.jpg",
	}
	assert.Equal(t, "/disk_mock/test/file.jpg", file.Url())
	assert.Equal(t, nil, file.Delete(context.Background()))
	file.Path = "/delete/error"
	assert.Equal(t, errors.New("Delete file error"), file.Delete(context.Background()))

	file.Disk = ""
	assert.Equal(t, "", file.Url())
	assert.Equal(t, errors.New("disk or path is empty"), file.Delete(context.Background()))

}

func TestPermission(t *testing.T) {
	permissionFilter := &entities.PermissionFilter{
		Filter: &entities.Filter{
			BaseUrl:         "/permission",
			Search:          "test",
			Page:            2,
			Limit:           10,
			Sorts:           []*entities.Sort{{"created_at", "desc"}},
			IgnoreUrlParams: []string{},
			ExcludeIDs:      []int{1, 2},
		},
		RoleIDs: []int{1, 2},
	}

	assert.Equal(t, "/permission?q=test&role=1", permissionFilter.Base())
}

func TestPost(t *testing.T) {
	post := &entities.Post{ID: 1, Slug: "test-post"}
	assert.Equal(t, config.Url("/test-post-1.html"), post.Url())
	postFilter := &entities.PostFilter{
		Filter: &entities.Filter{
			BaseUrl:         "/post",
			Search:          "test",
			Page:            2,
			Limit:           10,
			Sorts:           []*entities.Sort{{"created_at", "desc"}},
			IgnoreUrlParams: []string{},
			ExcludeIDs:      []int{1, 2},
		},
		Publish:  "draft",
		TopicIDs: []int{1, 2},
		UserIDs:  []int{1, 2},
	}

	assert.Equal(t, "/post?publish=draft&q=test&topic=1&user=1", postFilter.Base())
}

func TestRole(t *testing.T) {
	assert.Equal(t, "all", entities.PERM_ALL.String())
	assert.Equal(t, "own", entities.PERM_OWN.String())
	assert.Equal(t, "none", entities.PERM_NONE.String())

	var PERM_TEST entities.PermType = "test"
	assert.Equal(t, "none", PERM_TEST.String())

	assert.Equal(t, entities.PERM_ALL, entities.GetPermTypeValue("all"))
	assert.Equal(t, entities.PERM_OWN, entities.GetPermTypeValue("own"))
	assert.Equal(t, entities.PERM_NONE, entities.GetPermTypeValue("none"))
	assert.Equal(t, entities.PERM_NONE, entities.GetPermTypeValue("test"))

	roleFilter := &entities.RoleFilter{
		Filter: &entities.Filter{
			BaseUrl:         "/role",
			Search:          "test",
			Page:            2,
			Limit:           10,
			Sorts:           []*entities.Sort{{"created_at", "desc"}},
			IgnoreUrlParams: []string{},
			ExcludeIDs:      []int{1, 2},
		},
	}

	assert.Equal(t, "/role?q=test", roleFilter.Base())
}

func TestTopic(t *testing.T) {
	topic := &entities.Topic{ID: 1, Slug: "test-topic"}
	assert.Equal(t, config.Url("/test-topic"), topic.Url())
	assert.Equal(t, config.Url("/test-topic/feed"), topic.FeedUrl())

	topic1 := &entities.Topic{ID: 1, Name: "Topic 1", Slug: "test-topic-1"}
	topic2 := &entities.Topic{ID: 2, Name: "Topic 2", Slug: "test-topic-2"}
	topic3 := &entities.Topic{ID: 3, Name: "Topic 3", Slug: "test-topic-3"}
	topic4 := &entities.Topic{ID: 4, Name: "Topic 4", Slug: "test-topic-4"}
	topic5 := &entities.Topic{ID: 5, Name: "Topic 5", Slug: "test-topic-5"}

	topics := []*entities.Topic{topic1, topic2, topic3, topic4, topic5}
	assert.Equal(t, topics, entities.PrintTopicsTree(topics, []int{0}))
	assert.Equal(t, topics[:4], entities.PrintTopicsTree(topics, []int{5}))

	topic5.ParentID = 4
	topic5.Parent = topic4
	topic4.ParentID = 3
	topic4.Parent = topic3

	topic3.Children = []*entities.Topic{topic4}
	topic4.Children = []*entities.Topic{topic5}
	topicTreePrint := []*entities.Topic{
		topic1,
		topic2,
		topic3,
		topic4,
		topic5,
	}

	// assert.Equal(t, topicTreePrint, entities.GetTopicsTree(topics, 0, 0, []int{0}))

	print := entities.PrintTopicsTree(topicTreePrint, []int{0})
	assert.Equal(t, print[3].Name, "--Topic 4")
	assert.Equal(t, print[4].Name, "----Topic 5")

	assert.Equal(t, []*entities.Topic{
		topic1,
		topic2,
		topic3,
	}, entities.GetTopicsTree(topics, 0, 0, []int{4}))

	topicFilter := &entities.TopicFilter{
		Filter: &entities.Filter{
			BaseUrl:         "/topic",
			Search:          "test",
			Page:            2,
			Limit:           10,
			Sorts:           []*entities.Sort{{"created_at", "desc"}},
			IgnoreUrlParams: []string{},
			ExcludeIDs:      []int{1, 2},
		},
	}

	assert.Equal(t, "/topic?q=test", topicFilter.Base())
}

func TestUser(t *testing.T) {
	user := &entities.User{ID: 1, Username: "test-user"}
	assert.Equal(t, config.Url("/u/test-user"), user.Url())
	assert.Equal(t, "", user.Avatar())
	user.ProviderAvatar = "http://provider.local/avatar.png"
	assert.Equal(t, "http://provider.local/avatar.png", user.Avatar())
	user.AvatarImage = &entities.File{
		ID:   1,
		Disk: "disk_mock",
		Path: "test/file.jpg",
	}
	assert.Equal(t, "/disk_mock/test/file.jpg", user.Avatar())
	assert.Equal(t, false, user.IsRoot())
	user.Roles = []*entities.Role{{
		ID:   2,
		Name: "User",
		Root: false,
	}}
	assert.Equal(t, false, user.IsRoot())
	user.Roles = []*entities.Role{{
		ID:   1,
		Name: "Admin",
		Root: true,
	}}
	assert.Equal(t, true, user.IsRoot())
	assert.Equal(t, "test-user", user.Name())
	user.DisplayName = "Test User"
	assert.Equal(t, "Test User", user.Name())

	token, err := user.JwtClaim(time.Now().Add(time.Hour*24*7), map[string]interface{}{
		"test": "test",
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, true, token != "")

	userFilter := &entities.UserFilter{
		Filter: &entities.Filter{
			BaseUrl:         "/user",
			Search:          "test",
			Page:            2,
			Limit:           10,
			Sorts:           []*entities.Sort{{"created_at", "desc"}},
			IgnoreUrlParams: []string{},
			ExcludeIDs:      []int{1, 2},
		},
	}

	assert.Equal(t, "/user?q=test", userFilter.Base())
}
