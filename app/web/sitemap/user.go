package websitemap

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
)

func User(c server.Context) error {
	sitemapUser := SitemapUrlSets{
		Xmlns:               "http://www.sitemaps.org/schemas/sitemap/0.9",
		XmlnsXsi:            "http://www.w3.org/2001/XMLSchema-instance",
		XmlnsImage:          "http://www.google.com/schemas/sitemap-image/1.1",
		XmlnsSchemaLocation: "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd http://www.google.com/schemas/sitemap-image/1.1 http://www.google.com/schemas/sitemap-image/1.1/sitemap-image.xsd",
	}

	users, err := repositories.User.Find(c.Context(), &entities.UserFilter{Filter: &entities.Filter{Limit: 10000}})

	if err != nil {
		c.Logger().Error(err)
		return c.Status(http.StatusInternalServerError).SendString("Error")
	}

	for _, user := range users {
		avatar := user.Avatar()
		sitemapUrl := &SitemapUrl{
			Loc: &SitemapLoc{
				Value: user.Url(),
			},
			LastMod: &SitemapLastMod{
				Value: user.UpdatedAt.Format(time.RFC3339),
			},
		}

		if avatar != "" {
			sitemapUrl.Image = &SitemapImage{
				Loc: &SitemapImageLoc{
					Value: avatar,
				},
			}
		}
		sitemapUser.Urls = append(sitemapUser.Urls, sitemapUrl)
	}

	sitemapBytes, err := xml.Marshal(sitemapUser)

	if err != nil {
		c.Logger().Error(err)
		return c.Status(http.StatusInternalServerError).SendString("Error")
	}

	c.Response().Header("content-type", "text/xml; charset=UTF-8")
	return c.SendString(`<?xml version="1.0" encoding="UTF-8"?>` + string(sitemapBytes))
}
