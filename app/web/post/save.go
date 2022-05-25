package webpost

import (
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/services"
	"github.com/ngocphuongnb/tetua/app/utils"
)

func postMutationToPost(postData *entities.PostMutation) *entities.Post {
	return &entities.Post{
		Name:            postData.Name,
		Slug:            slug.Make(postData.Name),
		Content:         postData.Content,
		ContentHTML:     postData.ContentHTML,
		Description:     postData.Description,
		FeaturedImageID: postData.FeaturedImageID,
		Draft:           postData.Draft,
		TopicIDs:        postData.TopicIDs,
	}
}

func Save(c server.Context) (err error) {
	var post *entities.Post
	featuredImage := &entities.File{}
	postData := getPostSaveData(c)
	contentHtml, err := utils.MarkdownToHtml(postData.Content)

	if err != nil {
		c.WithError("Error convert markdown to html", err)
	}

	if post = c.Post(); post != nil {
		postData.Slug = post.Slug
		postData.Description = post.Description
	}

	if postData.FeaturedImageID > 0 {
		if featuredImage, err = repositories.File.ByID(c.Context(), postData.FeaturedImageID); err != nil {
			c.WithError("Error getting featured image", err)
		}
	}

	if !c.Messages().HasError() {
		var savedPost *entities.Post
		postData.Slug = slug.Make(postData.Name)
		postData.ContentHTML = contentHtml
		savePostData := postMutationToPost(postData)
		user := c.User()

		if post := c.Post(); post != nil {
			now := time.Now()
			savePostData.ID = post.ID
			savePostData.Approved = post.Approved
			savePostData.UpdatedAt = &now
			savedPost, err = repositories.Post.Update(c.Context(), savePostData)
		} else {
			savePostData.UserID = user.ID
			savePostData.Approved = user.IsRoot() || config.Setting("auto_approve_post") == "yes"
			savedPost, err = repositories.Post.Create(c.Context(), savePostData)
		}

		if err != nil {
			c.Logger().Error("Error creating post", err)
			c.Messages().AppendError("Error saving post")
			return getComposeView(c, postData, featuredImage)
		}

		return c.RedirectToRoute("post.compose", entities.Map{"id": savedPost.ID})
	}

	return getComposeView(c, postData, featuredImage)
}

func getPostSaveData(c server.Context) *entities.PostMutation {
	var err error
	postData := &entities.PostMutation{}
	if err := c.BodyParser(postData); err != nil {
		c.Logger().Error("Error parsing body", err)
		c.Messages().AppendError("Bad request")
		return postData
	}

	postData.Content = utils.SanitizeMarkdown(postData.Content)
	lines := strings.Split(postData.Content, "\n")
	postData.Name = utils.SanitizePlainText(strings.Trim(strings.Trim(lines[0], "#"), " "))
	postData.Content = strings.Join(lines[1:], "\n")

	if postData.Name, err = utils.MarkdownToHtml(postData.Name); err != nil {
		c.WithError("Error convert markdown to html", err)
		postData.Name = ""
	} else {
		postData.Name = utils.SanitizePlainText(postData.Name)
		postData.Name = strings.ReplaceAll(postData.Name, "\n", "")
		postData.Name = strings.ReplaceAll(postData.Name, "\r", "")
		postData.Name = strings.ReplaceAll(postData.Name, "\t", "")
	}

	if featuredImage, err := services.SaveFile(c, "featured_image"); err != nil {
		c.WithError("Error saving featured image", err)
	} else if featuredImage != nil {
		postData.FeaturedImageID = featuredImage.ID
	}

	if postData.Name == "" || len(postData.Name) > 250 {
		c.Messages().AppendError("Name is required and can't be more than 250 characters")
	}

	if postData.Content == "" {
		c.Messages().AppendError("Content is required")
	}

	if len(postData.TopicIDs) == 0 {
		c.Messages().AppendError("Topic is required")
	}

	return postData
}
