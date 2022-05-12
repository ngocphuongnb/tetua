package websitemap

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
)

func Post(c server.Context) error {
	currentPage := c.ParamInt("page", 1)
	if currentPage < 1 {
		return c.Status(http.StatusInternalServerError).SendString("Error")
	}

	posts, err := repositories.Post.Find(c.Context(), &entities.PostFilter{
		Filter: &entities.Filter{
			Page:  currentPage,
			Limit: SITEMAP_PAGESIZE,
			Sorts: []*entities.Sort{{
				Field: "id",
				Order: "ASC",
			}},
		},
	})

	if err != nil {
		c.Logger().Error(err)
		return c.Status(http.StatusInternalServerError).SendString("Error")
	}

	sitemapPost := SitemapUrlSets{
		Xmlns:               "http://www.sitemaps.org/schemas/sitemap/0.9",
		XmlnsXsi:            "http://www.w3.org/2001/XMLSchema-instance",
		XmlnsImage:          "http://www.google.com/schemas/sitemap-image/1.1",
		XmlnsSchemaLocation: "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd http://www.google.com/schemas/sitemap-image/1.1 http://www.google.com/schemas/sitemap-image/1.1/sitemap-image.xsd",
	}

	for _, post := range posts {
		sitemapUrl := &SitemapUrl{
			Loc: &SitemapLoc{
				Value: post.Url(),
			},
			LastMod: &SitemapLastMod{
				Value: post.UpdatedAt.Format(time.RFC3339),
			},
		}

		if post.FeaturedImage != nil && post.FeaturedImage.Url() != "" {
			sitemapUrl.Image = &SitemapImage{
				Loc: &SitemapImageLoc{
					Value: post.FeaturedImage.Url(),
				},
			}
		}
		sitemapPost.Urls = append(sitemapPost.Urls, sitemapUrl)
	}

	sitemapBytes, err := xml.Marshal(sitemapPost)

	if err != nil {
		c.Logger().Error(err)
		return c.Status(http.StatusInternalServerError).SendString("Error")
	}

	c.Response().Header("content-type", "text/xml; charset=UTF-8")
	return c.SendString(`<?xml version="1.0" encoding="UTF-8"?>` + string(sitemapBytes))
}
