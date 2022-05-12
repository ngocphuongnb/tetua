package entrepository

import (
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent"
)

func entFileToFile(file *ent.File) *entities.File {
	if file == nil {
		return nil
	}
	f := &entities.File{
		ID:        file.ID,
		CreatedAt: &file.CreatedAt,
		UpdatedAt: &file.UpdatedAt,
		DeletedAt: &file.DeletedAt,
		Disk:      file.Disk,
		Path:      file.Path,
		Type:      file.Type,
		Size:      file.Size,
		UserID:    file.UserID,
	}

	if file.Edges.Posts != nil {
		f.Posts = entPostsToPosts(file.Edges.Posts)
	}
	if file.Edges.User != nil {
		f.User = entUserToUser(file.Edges.User)
	}

	return f
}

func entFilesToFiles(files []*ent.File) []*entities.File {
	var result []*entities.File

	for _, file := range files {
		result = append(result, entFileToFile(file))
	}

	return result
}

func entPermissionToPermission(permission *ent.Permission) *entities.Permission {
	if permission == nil {
		return nil
	}
	p := &entities.Permission{
		ID:        permission.ID,
		RoleID:    permission.RoleID,
		Action:    permission.Action,
		Value:     permission.Value,
		CreatedAt: &permission.CreatedAt,
		UpdatedAt: &permission.UpdatedAt,
		DeletedAt: &permission.DeletedAt,
	}

	if permission.Edges.Role != nil {
		p.Role = entRoleToRole(permission.Edges.Role)
	}

	return p
}

func entPermissionsToPermissions(permissions []*ent.Permission) []*entities.Permission {
	var result []*entities.Permission

	for _, permission := range permissions {
		result = append(result, entPermissionToPermission(permission))
	}

	return result
}

func entPostToPost(post *ent.Post) *entities.Post {
	if post == nil {
		return nil
	}
	p := &entities.Post{
		ID:              post.ID,
		Name:            post.Name,
		Description:     post.Description,
		Slug:            post.Slug,
		Content:         post.Content,
		ContentHTML:     post.ContentHTML,
		Draft:           post.Draft,
		ViewCount:       post.ViewCount,
		CommentCount:    post.CommentCount,
		RatingCount:     post.RatingCount,
		RatingTotal:     post.RatingTotal,
		FeaturedImageID: post.FeaturedImageID,
		CreatedAt:       &post.CreatedAt,
		UpdatedAt:       &post.UpdatedAt,
		DeletedAt:       &post.DeletedAt,
		UserID:          post.UserID,
		Approved:        post.Approved,
		User:            &entities.User{},
		FeaturedImage:   &entities.File{},
		Topics:          entTopicsToTopics(post.Edges.Topics),
	}

	if post.Edges.FeaturedImage != nil {
		p.FeaturedImage = entFileToFile(post.Edges.FeaturedImage)
	}

	if post.Edges.User != nil {
		p.User = entUserToUser(post.Edges.User)
	}

	return p
}

func entPostsToPosts(posts []*ent.Post) []*entities.Post {
	var result []*entities.Post

	for _, post := range posts {
		result = append(result, entPostToPost(post))
	}

	return result
}

func entRoleToRole(role *ent.Role) *entities.Role {
	if role == nil {
		return nil
	}
	r := &entities.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Root:        role.Root,
		CreatedAt:   &role.CreatedAt,
		UpdatedAt:   &role.UpdatedAt,
		DeletedAt:   &role.DeletedAt,
	}

	if role.Edges.Users != nil {
		r.Users = entUsersToUsers(role.Edges.Users)
	}
	if role.Edges.Permissions != nil {
		r.Permissions = entPermissionsToPermissions(role.Edges.Permissions)
	}

	return r
}

func entRolesToRoles(roles []*ent.Role) []*entities.Role {
	var result []*entities.Role

	for _, role := range roles {
		result = append(result, entRoleToRole(role))
	}

	return result
}

func entTopicToTopic(topic *ent.Topic) *entities.Topic {
	if topic == nil {
		return nil
	}
	f := &entities.Topic{
		ID:          topic.ID,
		CreatedAt:   &topic.CreatedAt,
		UpdatedAt:   &topic.UpdatedAt,
		DeletedAt:   &topic.DeletedAt,
		Name:        topic.Name,
		Slug:        topic.Slug,
		Description: topic.Description,
		Content:     topic.Content,
		ContentHTML: topic.ContentHTML,
		ParentID:    topic.ParentID,
	}

	if topic.Edges.Parent != nil {
		f.Parent = entTopicToTopic(topic.Edges.Parent)
	}
	if topic.Edges.Children != nil {
		f.Children = entTopicsToTopics(topic.Edges.Children)
	}

	return f
}

func entTopicsToTopics(topics []*ent.Topic) []*entities.Topic {
	var result []*entities.Topic

	for _, topic := range topics {
		result = append(result, entTopicToTopic(topic))
	}

	return result
}

func entUserToUser(user *ent.User) *entities.User {
	if user == nil {
		return nil
	}
	u := &entities.User{
		ID:               user.ID,
		Username:         user.Username,
		Password:         user.Password,
		DisplayName:      user.DisplayName,
		URL:              user.URL,
		Provider:         user.Provider,
		ProviderID:       user.ProviderID,
		ProviderUsername: user.ProviderUsername,
		ProviderAvatar:   user.ProviderAvatar,
		Email:            user.Email,
		Bio:              user.Bio,
		BioHTML:          user.BioHTML,
		Roles:            []*entities.Role{},
		Active:           user.Active,
		AvatarImageID:    user.AvatarImageID,
		CreatedAt:        &user.CreatedAt,
		UpdatedAt:        &user.UpdatedAt,
		DeletedAt:        &user.DeletedAt,
	}

	if user.Edges.Roles != nil {
		u.Roles = entRolesToRoles(user.Edges.Roles)
	}

	if user.Edges.AvatarImage != nil {
		u.AvatarImage = entFileToFile(user.Edges.AvatarImage)
	}

	return u
}

func entUsersToUsers(users []*ent.User) []*entities.User {
	var result []*entities.User

	for _, user := range users {
		result = append(result, entUserToUser(user))
	}

	return result
}

func entCommentToComment(comment *ent.Comment) *entities.Comment {
	if comment == nil {
		return nil
	}
	f := &entities.Comment{
		ID:          comment.ID,
		CreatedAt:   &comment.CreatedAt,
		UpdatedAt:   &comment.UpdatedAt,
		DeletedAt:   &comment.DeletedAt,
		UserID:      comment.UserID,
		ParentID:    comment.ParentID,
		Content:     comment.Content,
		ContentHTML: comment.ContentHTML,
	}

	if comment.Edges.Post != nil {
		f.Post = entPostToPost(comment.Edges.Post)
	}

	if comment.Edges.Parent != nil {
		f.Parent = entCommentToComment(comment.Edges.Parent)
	}

	if comment.Edges.User != nil {
		f.User = entUserToUser(comment.Edges.User)
	}

	return f
}

func entCommentsToComments(comments []*ent.Comment) []*entities.Comment {
	var result []*entities.Comment

	for _, comment := range comments {
		result = append(result, entCommentToComment(comment))
	}

	return result
}

// func entSettingToSetting(setting *ent.Setting) *entities.Setting {
// 	if setting == nil {
// 		return nil
// 	}
// 	f := &entities.Setting{
// 		ID:        setting.ID,
// 		Name:      setting.Name,
// 		Value:     setting.Value,
// 		CreatedAt: &setting.CreatedAt,
// 		UpdatedAt: &setting.UpdatedAt,
// 		DeletedAt: &setting.DeletedAt,
// 	}

// 	return f
// }

// func entSettingsToSettings(settings []*ent.Setting) []*entities.Setting {
// 	var result []*entities.Setting

// 	for _, setting := range settings {
// 		result = append(result, entSettingToSetting(setting))
// 	}

// 	return result
// }

func getSortFNs(sorts []*entities.Sort) []ent.OrderFunc {
	var result []ent.OrderFunc

	for _, sort := range sorts {
		direction := ent.Desc

		if sort.Order == "ASC" {
			direction = ent.Asc
		}

		result = append(result, direction(sort.Field))
	}

	return result
}
