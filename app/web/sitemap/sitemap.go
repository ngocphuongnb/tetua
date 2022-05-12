package websitemap

import (
	"encoding/xml"
	"fmt"
	"math"
	"net/http"

	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
)

func Index(c server.Context) error {

	total, err := repositories.Post.Count(c.Context())

	if err != nil {
		c.Logger().Error(err)
		return c.Status(http.StatusInternalServerError).SendString("Error")
	}

	totalPages := int(math.Ceil(float64(total) / float64(SITEMAP_PAGESIZE)))
	sitemapIndex := SitemapItems{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	sitemapIndex.Items = []*SitemapItem{{
		Loc: &SitemapLoc{
			Value: config.Url("/sitemap/topics.xml"),
		},
	}, {
		Loc: &SitemapLoc{
			Value: config.Url("/sitemap/users.xml"),
		},
	}}

	for i := 0; i < totalPages; i++ {
		sitemapIndex.Items = append(sitemapIndex.Items, &SitemapItem{
			Loc: &SitemapLoc{
				Value: config.Url(fmt.Sprintf("/sitemap/posts-%d.xml", i+1)),
			},
		})
	}

	sitemapBytes, err := xml.Marshal(sitemapIndex)

	if err != nil {
		c.Logger().Error(err)
		return c.Status(http.StatusInternalServerError).SendString("Error")
	}

	c.Response().Header("content-type", "text/xml; charset=UTF-8")
	return c.SendString(`<?xml version="1.0" encoding="UTF-8"?>` + string(sitemapBytes))
}
