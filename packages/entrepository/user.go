package entrepository

import (
	"context"
	"fmt"
	"time"

	"github.com/ngocphuongnb/tetua/app/entities"
	e "github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent/user"
)

type UserRepository struct {
	*BaseRepository[e.User, ent.User, *ent.UserQuery, *e.UserFilter]
}

func userById(ctx context.Context, client *ent.Client, id int) (*ent.User, error) {
	return client.User.Query().
		Where(user.IDEQ(id)).
		WithRoles().
		WithAvatarImage().
		Only(ctx)
}

func (u *UserRepository) ByUsername(ctx context.Context, username string) (*entities.User, error) {
	user, err := u.Client.User.
		Query().
		Where(user.UsernameEQ(username)).
		WithRoles().
		WithAvatarImage().
		Only(ctx)
	if err != nil {
		return nil, EntError(err, fmt.Sprintf("user not found with username: %s", username))
	}

	return entUserToUser(user), nil
}

func (u *UserRepository) ByUsernameOrEmail(ctx context.Context, username, email string) ([]*entities.User, error) {
	user, err := u.Client.User.
		Query().
		Where(
			user.Or(
				user.UsernameEQ(username),
				user.EmailEQ(email),
			),
		).
		WithRoles().
		WithAvatarImage().
		All(ctx)
	if err != nil {
		return nil, EntError(err, fmt.Sprintf("user not found with username or email: %s %s", username, email))
	}

	return entUsersToUsers(user), nil
}

func (u *UserRepository) ByProvider(ctx context.Context, providerName, providerId string) (*entities.User, error) {
	user, err := u.Client.User.
		Query().
		Where(
			user.Provider(providerName),
			user.Provider(providerId),
		).
		WithRoles().
		WithAvatarImage().
		Only(ctx)

	if err != nil {
		return nil, EntError(err, fmt.Sprintf("user not found with provider: %s %s", providerName, providerId))
	}

	return entUserToUser(user), nil
}

func (ur *UserRepository) CreateIfNotExistsByProvider(ctx context.Context, userData *entities.User) (*entities.User, error) {
	u, err := ur.Client.User.
		Query().
		Where(
			user.Provider(userData.Provider),
			user.ProviderID(userData.ProviderID),
		).
		WithRoles().
		WithAvatarImage().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ur.Create(ctx, userData)
		}

		return nil, err
	}

	return entUserToUser(u), nil
}

func (ur *UserRepository) Setting(ctx context.Context, id int, userData *entities.SettingMutation) (*entities.User, error) {
	uu := ur.Client.User.UpdateOneID(id).
		SetUsername(userData.Username).
		SetDisplayName(userData.DisplayName).
		SetURL(userData.URL).
		SetBio(userData.Bio).
		SetBioHTML(userData.BioHTML).
		SetEmail(userData.Email)

	if userData.AvatarImageID > 0 {
		uu.SetAvatarImageID(userData.AvatarImageID)
	}

	if userData.Password != "" {
		uu.SetPassword(userData.Password)
	}

	user, err := uu.Save(ctx)

	if err != nil {
		return nil, err
	}

	return entUserToUser(user), nil
}

func CreateUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{
		BaseRepository: &BaseRepository[e.User, ent.User, *ent.UserQuery, *e.UserFilter]{
			Name:      "user",
			Client:    client,
			ConvertFn: entUserToUser,
			ByIDFn:    userById,
			DeleteByIDFn: func(ctx context.Context, client *ent.Client, id int) error {
				return client.User.DeleteOneID(id).Exec(ctx)
			},
			CreateFn: func(ctx context.Context, client *ent.Client, data *e.User) (*ent.User, error) {
				uc := client.User.Create().
					SetUsername(data.Username).
					SetDisplayName(data.DisplayName).
					SetURL(data.URL).
					SetBio(data.Bio).
					SetBioHTML(data.BioHTML).
					SetEmail(data.Email).
					SetProvider(data.Provider).
					SetProviderID(data.ProviderID).
					SetProviderUsername(data.ProviderUsername).
					SetProviderAvatar(data.ProviderAvatar).
					SetActive(data.Active)

				if data.AvatarImageID > 0 {
					uc.SetAvatarImageID(data.AvatarImageID)
				}

				if data.Provider == "local" {
					uc.SetProviderID(fmt.Sprintf("%d", time.Now().UnixMicro()))
				}

				if len(data.RoleIDs) > 0 {
					uc.AddRoleIDs(data.RoleIDs...)
				}

				if data.Password != "" {
					uc.SetPassword(data.Password)
				}

				user, err := uc.Save(ctx)

				if err != nil {
					return nil, err
				}

				if data.Provider == "local" {
					user, err = client.User.
						UpdateOneID(user.ID).
						SetProviderID(fmt.Sprintf("%d", user.ID)).Save(ctx)
					if err != nil {
						return nil, err
					}
				}

				return userById(ctx, client, user.ID)
			},
			UpdateFn: func(ctx context.Context, client *ent.Client, data *e.User) (*ent.User, error) {
				if data.ID == 0 {
					return nil, fmt.Errorf("user id is required")
				}
				uu := client.User.UpdateOneID(data.ID).
					SetUsername(data.Username).
					SetDisplayName(data.DisplayName).
					SetURL(data.URL).
					SetBio(data.Bio).
					SetBioHTML(data.BioHTML).
					SetEmail(data.Email).
					SetProvider(data.Provider).
					SetProviderID(data.ProviderID).
					SetProviderUsername(data.ProviderUsername).
					SetProviderAvatar(data.ProviderAvatar).
					SetActive(data.Active)

				if data.AvatarImageID > 0 {
					uu.SetAvatarImageID(data.AvatarImageID)
				}

				if len(data.RoleIDs) > 0 {
					oldUserEnt, err := userById(ctx, client, data.ID)
					if err != nil {
						return nil, err
					}
					oldUser := entUserToUser(oldUserEnt)
					oldRoleIDs := utils.SliceMap(oldUser.Roles, func(r *entities.Role) int {
						return r.ID
					})
					uu.RemoveRoleIDs(oldRoleIDs...)
					uu.AddRoleIDs(data.RoleIDs...)
				}

				if data.Password != "" {
					uu.SetPassword(data.Password)
				}

				user, err := uu.Save(ctx)

				if err != nil {
					return nil, err
				}

				return user, nil
			},
			QueryFilterFn: func(client *ent.Client, filters ...*e.UserFilter) *ent.UserQuery {
				query := client.User.Query().Where(user.DeletedAtIsNil())

				if len(filters) > 0 && filters[0].Search != "" {
					query = query.Where(user.UsernameContainsFold(filters[0].Search))
				}
				return query
			},
			FindFn: func(ctx context.Context, query *ent.UserQuery, filters ...*e.UserFilter) ([]*ent.User, error) {
				page, limit, sorts := getPaginateParams(filters...)
				return query.
					Limit(limit).
					Offset((page - 1) * limit).
					Order(sorts...).All(ctx)
			},
		},
	}
}
