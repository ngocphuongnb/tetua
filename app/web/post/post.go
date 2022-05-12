package webpost

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/views"
)

func List(c server.Context) error {
	paginate, err := repositories.Post.Paginate(c.Context(), &entities.PostFilter{
		Filter: &entities.Filter{
			BaseUrl:         config.Url("/posts"),
			Page:            c.QueryInt("page"),
			IgnoreUrlParams: []string{"user"},
		},
		UserIDs: []int{c.User().ID},
		Approve: "all",
		Publish: "all",
	})

	if err != nil {
		c.WithError("Something went wrong", err)
	}

	return c.Render(views.PostList(paginate))
}

func Delete(c server.Context) error {
	if err := repositories.Post.DeleteByID(c.Context(), c.ParamInt("id")); err != nil {
		c.Logger().Error("Error deleting post", err)
		return c.Status(http.StatusBadRequest).Json(&entities.Message{
			Type:    "error",
			Message: "Error deleting post",
		})
	}

	return c.Status(http.StatusOK).Json(&entities.Message{
		Type:    "success",
		Message: "Post deleted",
	})
}

func View(c server.Context) error {
	var post = &entities.Post{}
	var slug = c.Param("slug")
	var slugParts = strings.Split(slug, "-")

	if len(slugParts) == 0 {
		return c.Status(http.StatusNotFound).Render(views.Error("Post not found"))
	}

	var slugId = slugParts[len(slugParts)-1]
	var relatedPosts []*entities.Post
	var comments = []*entities.Comment{}
	var wg sync.WaitGroup
	var postId, err = strconv.Atoi(slugId)

	if err != nil {
		c.WithError("Invalid post id", err)
		return c.Status(http.StatusNotFound).Render(views.Error("Post not found"))
	}

	if post, err = repositories.Post.PublishedPostByID(c.Context(), postId); err != nil || post == nil {
		if err != nil {
			c.WithError("Error finding post", err)
		}
		return c.Status(http.StatusNotFound).Render(views.Error("Post not found"))
	}

	postSlug := fmt.Sprintf("%s-%d", post.Slug, post.ID)

	if postSlug != slug {
		return c.RedirectToRoute("post.view", entities.Map{"slug": postSlug})
	}

	wg.Add(2)
	go func() {
		if err := repositories.Post.IncreaseViewCount(c.Context(), postId, 1); err != nil {
			c.Logger().Error("Error incrementing post view", err)
		}
	}()

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		relatedPosts = getRelatedPosts(c, post)
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		comments, err = repositories.Comment.Find(c.Context(), &entities.CommentFilter{
			PostIDs: []int{postId},
			Filter: &entities.Filter{
				Limit: 10000,
			},
		})
	}(&wg)

	wg.Wait()
	c.Meta().Title = post.Name
	c.Meta().Description = post.Description

	if post.FeaturedImage != nil {
		c.Meta().Image = post.FeaturedImage.Url()
	}

	return c.Render(views.PostView(post, relatedPosts, comments))
}

func getRelatedPosts(c server.Context, post *entities.Post) []*entities.Post {
	var relatedPosts []*entities.Post
	var err error

	if relatedPosts, err = repositories.Post.Find(c.Context(), &entities.PostFilter{
		Filter: &entities.Filter{
			Limit:      8,
			ExcludeIDs: []int{post.ID},
			Sorts: []*entities.Sort{{
				Field: "view_count",
				Order: "desc",
			}},
		},
		TopicIDs: utils.SliceMap(post.Topics, func(topic *entities.Topic) int {
			return topic.ID
		}),
	}); err != nil {
		c.Logger().Error("Error finding related posts", err)
	}

	return relatedPosts
}
