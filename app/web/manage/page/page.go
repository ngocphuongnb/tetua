package managepage

import (
	"net/http"
	"time"

	"github.com/gosimple/slug"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/services"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/views"
)

func Index(c server.Context) error {
	c.Meta().Title = "Manage pages"
	page := c.QueryInt("page", 1)
	search := c.Query("q")
	publish := c.Query("publish", "all")
	status := http.StatusOK
	data, err := repositories.Page.Paginate(c.Context(), &entities.PageFilter{
		Publish: publish,
		Filter: &entities.Filter{
			BaseUrl: utils.Url("/manage/pages"),
			Page:    page,
			Search:  search,
		}})

	if err != nil {
		status = http.StatusBadRequest
		c.WithError("Error getting pages", err)
	}

	return c.Status(status).Render(views.ManagePageIndex(
		data,
		search,
		publish,
	))
}

func Delete(c server.Context) error {
	page, err := getProcessingPage(c)

	if err != nil {
		c.Logger().Error("Error deleting page", err)
		return c.Status(http.StatusBadRequest).SendString("Error deleting page")
	}

	if err := repositories.Page.DeleteByID(c.Context(), page.ID); err != nil {
		c.Logger().Error("Error deleting page", err)
		return c.Status(http.StatusBadRequest).SendString("Error deleting page")
	}

	return c.Status(http.StatusOK).SendString("Page deleted")
}

func Compose(c server.Context) (err error) {
	featuredImage := &entities.File{}
	page := &entities.Page{}
	c.Meta().Title = "Create Page"

	if page, err = getProcessingPage(c); err != nil {
		return err
	}

	if page.ID > 0 {
		c.Meta().Title = "Edit Page: " + page.Name
		featuredImage = page.FeaturedImage
	}

	return getComposeView(c, page, featuredImage)
}

func Save(c server.Context) (err error) {
	var page *entities.Page
	featuredImage := &entities.File{}
	pageData := getPageSaveData(c)
	contentHtml, err := utils.MarkdownToHtml(pageData.Content)

	if err != nil {
		c.WithError("Error convert markdown to html", err)
	}

	if page, err = getProcessingPage(c); err != nil {
		return err
	}

	if pageData.FeaturedImageID > 0 {
		if featuredImage, err = repositories.File.ByID(c.Context(), pageData.FeaturedImageID); err != nil {
			c.WithError("Error getting featured image", err)
		}
	}

	if !c.Messages().HasError() {
		var savedPage *entities.Page
		pageData.ContentHTML = contentHtml

		if page.ID > 0 {
			now := time.Now()
			pageData.ID = page.ID
			pageData.UpdatedAt = &now
			savedPage, err = repositories.Page.Update(c.Context(), pageData)
		} else {
			savedPage, err = repositories.Page.Create(c.Context(), pageData)
		}

		if err != nil {
			c.WithError("Error saving page", err)
			return getComposeView(c, pageData, featuredImage)
		}

		return c.RedirectToRoute("manage.page.compose", entities.Map{"id": savedPage.ID})
	}

	return getComposeView(c, pageData, featuredImage)
}

func getComposeView(c server.Context, data *entities.Page, featuredImage *entities.File) (err error) {
	status := http.StatusOK

	if err != nil {
		c.Logger().Error("Error getting pages", err)
		c.Messages().AppendError("Error getting pages")
	}

	if c.Messages().HasError() {
		status = http.StatusBadRequest
	}

	return c.Status(status).Render(views.ManagePageCompose(data, featuredImage))
}

func getPageSaveData(c server.Context) *entities.Page {
	pageData := &entities.Page{}

	if err := c.BodyParser(pageData); err != nil {
		c.WithError("Bad request", err)
		return pageData
	}

	pageData.Content = utils.SanitizeMarkdown(pageData.Content)
	pageData.Name = utils.SanitizePlainText(pageData.Name)

	if pageData.Slug == "" {
		pageData.Slug = slug.Make(pageData.Name)
	} else {
		pageData.Slug = utils.SanitizePlainText(pageData.Slug)
	}

	if featuredImage, err := services.SaveFile(c, "featured_image"); err != nil {
		c.WithError("Error saving featured image", err)
	} else if featuredImage != nil {
		pageData.FeaturedImageID = featuredImage.ID
	}

	if pageData.Name == "" || len(pageData.Name) > 250 {
		c.Messages().AppendError("Name is required and can't be more than 250 characters")
	}

	if pageData.Content == "" {
		c.Messages().AppendError("Content is required")
	}

	return pageData
}

func getProcessingPage(c server.Context) (page *entities.Page, err error) {
	if c.Param("id") == "new" {
		return &entities.Page{}, nil
	}
	return repositories.Page.ByID(c.Context(), c.ParamInt("id"))
}
