package web

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/feeds"
	"github.com/ngocphuongnb/tetua/app/cache"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/views"
)

func TopicView(c server.Context) (err error) {
	topicSlug := c.Param("slug")
	topics := utils.SliceFilter(cache.Topics, func(t *entities.Topic) bool {
		return t.Slug == topicSlug
	})

	if len(topics) == 0 {
		c.Meta().Title = "Topic not found"
		return c.Status(http.StatusNotFound).Render(views.Error("Topic not found"))
	}

	topic := topics[0]
	c.Meta().Title = topic.Name
	var wg sync.WaitGroup
	var paginate *entities.Paginate[entities.Post]
	var topPosts []*entities.Post

	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		paginate, err = repositories.Post.Paginate(c.Context(), &entities.PostFilter{
			Filter: &entities.Filter{
				Page:            c.QueryInt("page"),
				IgnoreUrlParams: []string{"topic"},
			},
			TopicIDs: []int{topic.ID},
		})
	}(&wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var topPostErr error
		topPosts, topPostErr = repositories.Post.Find(c.Context(), &entities.PostFilter{
			Filter: &entities.Filter{
				Limit: 8,
				Sorts: []*entities.Sort{{
					Field: "view_count",
					Order: "desc",
				}},
			},
			TopicIDs: []int{topic.ID},
		})

		if topPostErr != nil {
			c.Logger().Error(topPostErr)
		}
	}(&wg)

	wg.Wait()

	if err != nil {
		c.Logger().Error(err)
		return c.Status(http.StatusBadGateway).Render(views.Error("Something went wrong"))
	}

	return c.Render(views.TopicView(cache.Topics, topic, paginate, topPosts))
}

func TopicFeed(c server.Context) error {
	topics := utils.SliceFilter(cache.Topics, func(t *entities.Topic) bool {
		return t.Slug == c.Param("slug")
	})

	if len(topics) == 0 {
		c.Meta().Title = "Topic not found"
		return c.Status(http.StatusNotFound).SendString("Topic not found")
	}

	topic := topics[0]
	posts, err := repositories.Post.Find(c.Context(), &entities.PostFilter{
		Filter: &entities.Filter{
			Limit: 50,
		},
		TopicIDs: []int{topic.ID},
	})

	if err != nil {
		c.Logger().Error(err)
		return c.Status(http.StatusNotFound).SendString("Topic not found")
	}

	feed := &feeds.Feed{
		Title:       topic.Name,
		Link:        &feeds.Link{Href: topic.Url()},
		Description: topic.Description,
		Author:      &feeds.Author{Name: config.Setting("contact_name"), Email: config.Setting("contact_email")},
	}

	feed.Items = utils.SliceMap(posts, func(post *entities.Post) *feeds.Item {
		return &feeds.Item{
			Id:          strconv.Itoa(post.ID),
			Title:       post.Name,
			Description: post.Description,
			Content:     post.ContentHTML,
			Author:      &feeds.Author{Name: post.User.Name(), Email: post.User.Email},
			Link:        &feeds.Link{Href: post.Url()},
			Created:     *post.CreatedAt,
			Updated:     *post.UpdatedAt,
		}
	})

	rss, err := feed.ToRss()
	if err != nil {
		c.Logger().Error(err)
		return c.Status(http.StatusInternalServerError).SendString("Error")
	}

	c.Response().Header("content-type", "application/xml; charset=utf-8")
	return c.SendString(rss)
}
