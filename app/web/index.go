package web

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/ngocphuongnb/tetua/app/cache"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/views"
)

func Index(c server.Context) (err error) {
	var page = c.QueryInt("page")
	var wg sync.WaitGroup
	var paginate *entities.Paginate[entities.Post]
	var topPosts []*entities.Post

	wg.Add(2)
	go func(wg *sync.WaitGroup, page int) {
		defer wg.Done()
		paginate, err = repositories.Post.Paginate(c.Context(), &entities.PostFilter{
			Filter: &entities.Filter{Page: page},
		})
	}(&wg, page)
	go func(wg *sync.WaitGroup, page int) {
		defer wg.Done()
		var topPostErr error
		topPosts, topPostErr = repositories.Post.Find(c.Context(), &entities.PostFilter{
			Filter: &entities.Filter{
				Page:  page,
				Limit: 8,
				Sorts: []*entities.Sort{{
					Field: "view_count",
					Order: "desc",
				}},
			}})
		if topPostErr != nil {
			c.Logger().Error(topPostErr)
		}
	}(&wg, page)

	wg.Wait()

	if err != nil {
		c.Logger().Error(err)
		return c.Status(http.StatusBadGateway).Render(views.Error("Something went wrong"))
	}

	return c.Render(views.Index(cache.Topics, paginate, topPosts))
}

func Search(c server.Context) (err error) {
	c.Meta().Title = "Search"
	var paginate *entities.Paginate[entities.Post]
	var searchQuery = c.Query("q")

	if searchQuery != "" {
		c.Meta().Query = searchQuery
		c.Meta().Title = searchQuery + " - Search result for " + searchQuery
		c.Meta().Canonical = config.Url(c.Path() + "?q=" + url.QueryEscape(searchQuery))
	}

	paginate, err = repositories.Post.Paginate(c.Context(), &entities.PostFilter{
		Filter: &entities.Filter{
			Page:   c.QueryInt("page"),
			Search: searchQuery,
		}})

	if err != nil {
		c.Logger().Error(err)
		return c.Status(http.StatusBadGateway).Render(views.Error("Something went wrong"))
	}

	return c.Render(views.Search(cache.Topics, paginate))
}
