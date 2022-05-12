package web

import (
	"net/http"
	"strconv"

	"github.com/gorilla/feeds"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
)

func Feed(c server.Context) error {
	posts, err := repositories.Post.Find(c.Context(), &entities.PostFilter{Filter: &entities.Filter{
		Limit: 50,
	}})

	if err != nil {
		c.Logger().Error(err)
		return c.Status(http.StatusInternalServerError).SendString("Error")
	}

	feed := &feeds.Feed{
		Title:       config.Setting("app_name"),
		Link:        &feeds.Link{Href: config.Url("")},
		Description: config.Setting("app_desc"),
		Author:      &feeds.Author{Name: config.Setting("contact_name"), Email: config.Setting("contact_email")},
		// Created:     time.Now(),
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
