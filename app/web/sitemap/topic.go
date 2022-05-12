package websitemap

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/ngocphuongnb/tetua/app/cache"
	"github.com/ngocphuongnb/tetua/app/server"
)

func Topic(c server.Context) error {
	sitemapTopic := SitemapUrlSets{
		Xmlns:               "http://www.sitemaps.org/schemas/sitemap/0.9",
		XmlnsXsi:            "http://www.w3.org/2001/XMLSchema-instance",
		XmlnsImage:          "http://www.google.com/schemas/sitemap-image/1.1",
		XmlnsSchemaLocation: "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd http://www.google.com/schemas/sitemap-image/1.1 http://www.google.com/schemas/sitemap-image/1.1/sitemap-image.xsd",
	}

	for _, topic := range cache.Topics {
		sitemapTopic.Urls = append(sitemapTopic.Urls, &SitemapUrl{
			Loc: &SitemapLoc{
				Value: topic.Url(),
			},
			LastMod: &SitemapLastMod{
				Value: topic.UpdatedAt.Format(time.RFC3339),
			},
		})
	}

	sitemapBytes, err := xml.Marshal(sitemapTopic)

	if err != nil {
		c.Logger().Error(err)
		return c.Status(http.StatusInternalServerError).SendString("Error")
	}

	c.Response().Header("content-type", "text/xml; charset=UTF-8")
	return c.SendString(`<?xml version="1.0" encoding="UTF-8"?>` + string(sitemapBytes))
}
