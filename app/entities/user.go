package entities

import (
	"fmt"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/utils"
)

// User is the model entity for the User schema.
type User struct {
	ID               int        `json:"id,omitempty" form:"id"`
	CreatedAt        *time.Time `json:"created_at,omitempty" form:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty" form:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" form:"deleted_at"`
	Username         string     `json:"username,omitempty" form:"username"`
	DisplayName      string     `json:"display_name,omitempty" form:"display_name"`
	URL              string     `json:"url,omitempty" form:"url"`
	Provider         string     `json:"provider,omitempty" form:"provider"`
	ProviderID       string     `json:"provider_id,omitempty" form:"provider_id"`
	ProviderUsername string     `json:"provider_username,omitempty" form:"provider_username"`
	ProviderAvatar   string     `json:"provider_avatar,omitempty" form:"provider_avatar"`
	Email            string     `json:"email,omitempty" form:"email"`
	Password         string     `json:"password,omitempty" form:"password"`
	Bio              string     `json:"bio,omitempty" form:"bio"`
	BioHTML          string     `json:"bio_html,omitempty" form:"bio_html"`
	RoleIDs          []int      `json:"role_ids,omitempty" form:"role_ids"`
	Roles            []*Role    `json:"roles,omitempty" form:"roles"`
	Active           bool       `json:"active,omitempty" form:"active"`
	AvatarImage      *File      `json:"avatar_image,omitempty" form:"avatar_image"`
	AvatarImageID    int        `json:"avatar_image_id,omitempty" form:"avatar_image_id"`
	AvatarImageUrl   string     `json:"avatar_image_url,omitempty" form:"avatar_image_url"`
}

type UserMutation struct {
	Username         string `json:"username,omitempty" form:"username"`
	DisplayName      string `json:"display_name,omitempty" form:"display_name"`
	URL              string `json:"url,omitempty" form:"url"`
	Provider         string `json:"provider,omitempty" form:"provider"`
	ProviderID       string `json:"provider_id,omitempty" form:"provider_id"`
	ProviderUsername string `json:"provider_username,omitempty" form:"provider_username"`
	ProviderAvatar   string `json:"provider_avatar,omitempty" form:"provider_avatar"`
	Email            string `json:"email,omitempty" form:"email"`
	Password         string `json:"password,omitempty" form:"password"`
	Bio              string `json:"bio,omitempty" form:"bio"`
	RoleIDs          []int  `json:"role_ids,omitempty" form:"role_ids"`
	Active           bool   `json:"active,omitempty" form:"active"`
}

type SettingMutation struct {
	Username      string `json:"username,omitempty" form:"username"`
	DisplayName   string `json:"display_name,omitempty" form:"display_name"`
	URL           string `json:"url,omitempty" form:"url"`
	Email         string `json:"email,omitempty" form:"email"`
	Password      string `json:"password,omitempty" form:"password"`
	Bio           string `json:"bio,omitempty" form:"bio"`
	BioHTML       string `json:"bio_html,omitempty" form:"bio_html"`
	AvatarImageID int    `json:"avatar_image_id,omitempty" form:"avatar_image_id"`
}

type UserJwtClaims struct {
	jwt.RegisteredClaims
	User User `json:"user"`
}

type UserFilter struct {
	*Filter
}

func (u *User) JwtClaim(exp time.Time, jwtHeaders ...map[string]interface{}) (string, error) {
	u.RoleIDs = make([]int, 0)

	for _, role := range u.Roles {
		if role.ID > 0 {
			u.RoleIDs = append(u.RoleIDs, role.ID)
		}
	}

	user := User{
		ID:             u.ID,
		Provider:       u.Provider,
		Username:       u.Username,
		Email:          u.Email,
		DisplayName:    u.DisplayName,
		Active:         u.Active,
		RoleIDs:        u.RoleIDs,
		AvatarImageUrl: u.Avatar(),
	}

	claims := &UserJwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Setting("app_name"),
			ExpiresAt: &jwt.NumericDate{Time: exp},
		},
		User: user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if len(jwtHeaders) > 0 {
		for k, v := range jwtHeaders[0] {
			token.Header[k] = v
		}
	}

	return token.SignedString([]byte(config.APP_KEY))
}

func (u *User) IsRoot() bool {
	if u == nil {
		return false
	}
	if u.Roles == nil {
		return false
	}
	for _, role := range u.Roles {
		if role.Root {
			return true
		}
	}

	return false
}

func (u *User) Name() string {
	if u == nil {
		return ""
	}
	if u.DisplayName != "" {
		return u.DisplayName
	}

	return u.Username
}
func (u *User) Url() string {
	if u == nil {
		return ""
	}
	return utils.Url("/u/" + u.Username)
}

func (u *User) Avatar() string {
	if u == nil {
		return ""
	}
	if u.AvatarImage != nil && u.AvatarImage.ID > 0 {
		return u.AvatarImage.Url()
	}

	if u.ProviderAvatar != "" {
		return u.ProviderAvatar
	}

	return u.AvatarImageUrl
}

func (u *User) AvatarElm(width, height string, disableLink bool) string {
	var userAvatar = u.Avatar()
	if userAvatar != "" {
		if !disableLink {
			return fmt.Sprintf(
				`<a class="avatar" href="%s" title="%s" target="_blank"><img src="%s" width="%s" height="%s" alt="%s" /></a>`,
				u.Url(),
				u.Name(),
				userAvatar,
				width,
				height,
				u.Name(),
			)
		} else {
			return fmt.Sprintf(
				`<img src="%s" width="%s" height="%s" alt="%s" />`,
				userAvatar,
				width,
				height,
				u.Name(),
			)
		}
	}

	return `<span class="avatar none"></span>`
}

func (p *UserFilter) Base() string {
	q := url.Values{}
	if p.Search != "" {
		q.Add("q", p.Search)
	}
	if queryString := q.Encode(); queryString != "" {
		return p.FilterBaseUrl() + "?" + q.Encode()
	}

	return p.FilterBaseUrl()
}
