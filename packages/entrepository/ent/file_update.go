// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent/file"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent/page"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent/post"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent/predicate"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent/user"
)

// FileUpdate is the builder for updating File entities.
type FileUpdate struct {
	config
	hooks    []Hook
	mutation *FileMutation
}

// Where appends a list predicates to the FileUpdate builder.
func (fu *FileUpdate) Where(ps ...predicate.File) *FileUpdate {
	fu.mutation.Where(ps...)
	return fu
}

// SetUpdatedAt sets the "updated_at" field.
func (fu *FileUpdate) SetUpdatedAt(t time.Time) *FileUpdate {
	fu.mutation.SetUpdatedAt(t)
	return fu
}

// SetDeletedAt sets the "deleted_at" field.
func (fu *FileUpdate) SetDeletedAt(t time.Time) *FileUpdate {
	fu.mutation.SetDeletedAt(t)
	return fu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (fu *FileUpdate) SetNillableDeletedAt(t *time.Time) *FileUpdate {
	if t != nil {
		fu.SetDeletedAt(*t)
	}
	return fu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (fu *FileUpdate) ClearDeletedAt() *FileUpdate {
	fu.mutation.ClearDeletedAt()
	return fu
}

// SetDisk sets the "disk" field.
func (fu *FileUpdate) SetDisk(s string) *FileUpdate {
	fu.mutation.SetDisk(s)
	return fu
}

// SetPath sets the "path" field.
func (fu *FileUpdate) SetPath(s string) *FileUpdate {
	fu.mutation.SetPath(s)
	return fu
}

// SetType sets the "type" field.
func (fu *FileUpdate) SetType(s string) *FileUpdate {
	fu.mutation.SetType(s)
	return fu
}

// SetSize sets the "size" field.
func (fu *FileUpdate) SetSize(i int) *FileUpdate {
	fu.mutation.ResetSize()
	fu.mutation.SetSize(i)
	return fu
}

// AddSize adds i to the "size" field.
func (fu *FileUpdate) AddSize(i int) *FileUpdate {
	fu.mutation.AddSize(i)
	return fu
}

// SetUserID sets the "user_id" field.
func (fu *FileUpdate) SetUserID(i int) *FileUpdate {
	fu.mutation.SetUserID(i)
	return fu
}

// SetNillableUserID sets the "user_id" field if the given value is not nil.
func (fu *FileUpdate) SetNillableUserID(i *int) *FileUpdate {
	if i != nil {
		fu.SetUserID(*i)
	}
	return fu
}

// ClearUserID clears the value of the "user_id" field.
func (fu *FileUpdate) ClearUserID() *FileUpdate {
	fu.mutation.ClearUserID()
	return fu
}

// SetUser sets the "user" edge to the User entity.
func (fu *FileUpdate) SetUser(u *User) *FileUpdate {
	return fu.SetUserID(u.ID)
}

// AddPostIDs adds the "posts" edge to the Post entity by IDs.
func (fu *FileUpdate) AddPostIDs(ids ...int) *FileUpdate {
	fu.mutation.AddPostIDs(ids...)
	return fu
}

// AddPosts adds the "posts" edges to the Post entity.
func (fu *FileUpdate) AddPosts(p ...*Post) *FileUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fu.AddPostIDs(ids...)
}

// AddPageIDs adds the "pages" edge to the Page entity by IDs.
func (fu *FileUpdate) AddPageIDs(ids ...int) *FileUpdate {
	fu.mutation.AddPageIDs(ids...)
	return fu
}

// AddPages adds the "pages" edges to the Page entity.
func (fu *FileUpdate) AddPages(p ...*Page) *FileUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fu.AddPageIDs(ids...)
}

// AddUserAvatarIDs adds the "user_avatars" edge to the User entity by IDs.
func (fu *FileUpdate) AddUserAvatarIDs(ids ...int) *FileUpdate {
	fu.mutation.AddUserAvatarIDs(ids...)
	return fu
}

// AddUserAvatars adds the "user_avatars" edges to the User entity.
func (fu *FileUpdate) AddUserAvatars(u ...*User) *FileUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return fu.AddUserAvatarIDs(ids...)
}

// Mutation returns the FileMutation object of the builder.
func (fu *FileUpdate) Mutation() *FileMutation {
	return fu.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (fu *FileUpdate) ClearUser() *FileUpdate {
	fu.mutation.ClearUser()
	return fu
}

// ClearPosts clears all "posts" edges to the Post entity.
func (fu *FileUpdate) ClearPosts() *FileUpdate {
	fu.mutation.ClearPosts()
	return fu
}

// RemovePostIDs removes the "posts" edge to Post entities by IDs.
func (fu *FileUpdate) RemovePostIDs(ids ...int) *FileUpdate {
	fu.mutation.RemovePostIDs(ids...)
	return fu
}

// RemovePosts removes "posts" edges to Post entities.
func (fu *FileUpdate) RemovePosts(p ...*Post) *FileUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fu.RemovePostIDs(ids...)
}

// ClearPages clears all "pages" edges to the Page entity.
func (fu *FileUpdate) ClearPages() *FileUpdate {
	fu.mutation.ClearPages()
	return fu
}

// RemovePageIDs removes the "pages" edge to Page entities by IDs.
func (fu *FileUpdate) RemovePageIDs(ids ...int) *FileUpdate {
	fu.mutation.RemovePageIDs(ids...)
	return fu
}

// RemovePages removes "pages" edges to Page entities.
func (fu *FileUpdate) RemovePages(p ...*Page) *FileUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fu.RemovePageIDs(ids...)
}

// ClearUserAvatars clears all "user_avatars" edges to the User entity.
func (fu *FileUpdate) ClearUserAvatars() *FileUpdate {
	fu.mutation.ClearUserAvatars()
	return fu
}

// RemoveUserAvatarIDs removes the "user_avatars" edge to User entities by IDs.
func (fu *FileUpdate) RemoveUserAvatarIDs(ids ...int) *FileUpdate {
	fu.mutation.RemoveUserAvatarIDs(ids...)
	return fu
}

// RemoveUserAvatars removes "user_avatars" edges to User entities.
func (fu *FileUpdate) RemoveUserAvatars(u ...*User) *FileUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return fu.RemoveUserAvatarIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (fu *FileUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	fu.defaults()
	if len(fu.hooks) == 0 {
		affected, err = fu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FileMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fu.mutation = mutation
			affected, err = fu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(fu.hooks) - 1; i >= 0; i-- {
			if fu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = fu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (fu *FileUpdate) SaveX(ctx context.Context) int {
	affected, err := fu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fu *FileUpdate) Exec(ctx context.Context) error {
	_, err := fu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fu *FileUpdate) ExecX(ctx context.Context) {
	if err := fu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (fu *FileUpdate) defaults() {
	if _, ok := fu.mutation.UpdatedAt(); !ok {
		v := file.UpdateDefaultUpdatedAt()
		fu.mutation.SetUpdatedAt(v)
	}
}

func (fu *FileUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   file.Table,
			Columns: file.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: file.FieldID,
			},
		},
	}
	if ps := fu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldUpdatedAt,
		})
	}
	if value, ok := fu.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldDeletedAt,
		})
	}
	if fu.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: file.FieldDeletedAt,
		})
	}
	if value, ok := fu.mutation.Disk(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldDisk,
		})
	}
	if value, ok := fu.mutation.Path(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldPath,
		})
	}
	if value, ok := fu.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldType,
		})
	}
	if value, ok := fu.mutation.Size(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: file.FieldSize,
		})
	}
	if value, ok := fu.mutation.AddedSize(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: file.FieldSize,
		})
	}
	if fu.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   file.UserTable,
			Columns: []string{file.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   file.UserTable,
			Columns: []string{file.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fu.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.PostsTable,
			Columns: []string{file.PostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: post.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.RemovedPostsIDs(); len(nodes) > 0 && !fu.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.PostsTable,
			Columns: []string{file.PostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: post.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.PostsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.PostsTable,
			Columns: []string{file.PostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: post.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fu.mutation.PagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.PagesTable,
			Columns: []string{file.PagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: page.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.RemovedPagesIDs(); len(nodes) > 0 && !fu.mutation.PagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.PagesTable,
			Columns: []string{file.PagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: page.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.PagesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.PagesTable,
			Columns: []string{file.PagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: page.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fu.mutation.UserAvatarsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.UserAvatarsTable,
			Columns: []string{file.UserAvatarsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.RemovedUserAvatarsIDs(); len(nodes) > 0 && !fu.mutation.UserAvatarsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.UserAvatarsTable,
			Columns: []string{file.UserAvatarsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fu.mutation.UserAvatarsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.UserAvatarsTable,
			Columns: []string{file.UserAvatarsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{file.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// FileUpdateOne is the builder for updating a single File entity.
type FileUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *FileMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (fuo *FileUpdateOne) SetUpdatedAt(t time.Time) *FileUpdateOne {
	fuo.mutation.SetUpdatedAt(t)
	return fuo
}

// SetDeletedAt sets the "deleted_at" field.
func (fuo *FileUpdateOne) SetDeletedAt(t time.Time) *FileUpdateOne {
	fuo.mutation.SetDeletedAt(t)
	return fuo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableDeletedAt(t *time.Time) *FileUpdateOne {
	if t != nil {
		fuo.SetDeletedAt(*t)
	}
	return fuo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (fuo *FileUpdateOne) ClearDeletedAt() *FileUpdateOne {
	fuo.mutation.ClearDeletedAt()
	return fuo
}

// SetDisk sets the "disk" field.
func (fuo *FileUpdateOne) SetDisk(s string) *FileUpdateOne {
	fuo.mutation.SetDisk(s)
	return fuo
}

// SetPath sets the "path" field.
func (fuo *FileUpdateOne) SetPath(s string) *FileUpdateOne {
	fuo.mutation.SetPath(s)
	return fuo
}

// SetType sets the "type" field.
func (fuo *FileUpdateOne) SetType(s string) *FileUpdateOne {
	fuo.mutation.SetType(s)
	return fuo
}

// SetSize sets the "size" field.
func (fuo *FileUpdateOne) SetSize(i int) *FileUpdateOne {
	fuo.mutation.ResetSize()
	fuo.mutation.SetSize(i)
	return fuo
}

// AddSize adds i to the "size" field.
func (fuo *FileUpdateOne) AddSize(i int) *FileUpdateOne {
	fuo.mutation.AddSize(i)
	return fuo
}

// SetUserID sets the "user_id" field.
func (fuo *FileUpdateOne) SetUserID(i int) *FileUpdateOne {
	fuo.mutation.SetUserID(i)
	return fuo
}

// SetNillableUserID sets the "user_id" field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableUserID(i *int) *FileUpdateOne {
	if i != nil {
		fuo.SetUserID(*i)
	}
	return fuo
}

// ClearUserID clears the value of the "user_id" field.
func (fuo *FileUpdateOne) ClearUserID() *FileUpdateOne {
	fuo.mutation.ClearUserID()
	return fuo
}

// SetUser sets the "user" edge to the User entity.
func (fuo *FileUpdateOne) SetUser(u *User) *FileUpdateOne {
	return fuo.SetUserID(u.ID)
}

// AddPostIDs adds the "posts" edge to the Post entity by IDs.
func (fuo *FileUpdateOne) AddPostIDs(ids ...int) *FileUpdateOne {
	fuo.mutation.AddPostIDs(ids...)
	return fuo
}

// AddPosts adds the "posts" edges to the Post entity.
func (fuo *FileUpdateOne) AddPosts(p ...*Post) *FileUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fuo.AddPostIDs(ids...)
}

// AddPageIDs adds the "pages" edge to the Page entity by IDs.
func (fuo *FileUpdateOne) AddPageIDs(ids ...int) *FileUpdateOne {
	fuo.mutation.AddPageIDs(ids...)
	return fuo
}

// AddPages adds the "pages" edges to the Page entity.
func (fuo *FileUpdateOne) AddPages(p ...*Page) *FileUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fuo.AddPageIDs(ids...)
}

// AddUserAvatarIDs adds the "user_avatars" edge to the User entity by IDs.
func (fuo *FileUpdateOne) AddUserAvatarIDs(ids ...int) *FileUpdateOne {
	fuo.mutation.AddUserAvatarIDs(ids...)
	return fuo
}

// AddUserAvatars adds the "user_avatars" edges to the User entity.
func (fuo *FileUpdateOne) AddUserAvatars(u ...*User) *FileUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return fuo.AddUserAvatarIDs(ids...)
}

// Mutation returns the FileMutation object of the builder.
func (fuo *FileUpdateOne) Mutation() *FileMutation {
	return fuo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (fuo *FileUpdateOne) ClearUser() *FileUpdateOne {
	fuo.mutation.ClearUser()
	return fuo
}

// ClearPosts clears all "posts" edges to the Post entity.
func (fuo *FileUpdateOne) ClearPosts() *FileUpdateOne {
	fuo.mutation.ClearPosts()
	return fuo
}

// RemovePostIDs removes the "posts" edge to Post entities by IDs.
func (fuo *FileUpdateOne) RemovePostIDs(ids ...int) *FileUpdateOne {
	fuo.mutation.RemovePostIDs(ids...)
	return fuo
}

// RemovePosts removes "posts" edges to Post entities.
func (fuo *FileUpdateOne) RemovePosts(p ...*Post) *FileUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fuo.RemovePostIDs(ids...)
}

// ClearPages clears all "pages" edges to the Page entity.
func (fuo *FileUpdateOne) ClearPages() *FileUpdateOne {
	fuo.mutation.ClearPages()
	return fuo
}

// RemovePageIDs removes the "pages" edge to Page entities by IDs.
func (fuo *FileUpdateOne) RemovePageIDs(ids ...int) *FileUpdateOne {
	fuo.mutation.RemovePageIDs(ids...)
	return fuo
}

// RemovePages removes "pages" edges to Page entities.
func (fuo *FileUpdateOne) RemovePages(p ...*Page) *FileUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fuo.RemovePageIDs(ids...)
}

// ClearUserAvatars clears all "user_avatars" edges to the User entity.
func (fuo *FileUpdateOne) ClearUserAvatars() *FileUpdateOne {
	fuo.mutation.ClearUserAvatars()
	return fuo
}

// RemoveUserAvatarIDs removes the "user_avatars" edge to User entities by IDs.
func (fuo *FileUpdateOne) RemoveUserAvatarIDs(ids ...int) *FileUpdateOne {
	fuo.mutation.RemoveUserAvatarIDs(ids...)
	return fuo
}

// RemoveUserAvatars removes "user_avatars" edges to User entities.
func (fuo *FileUpdateOne) RemoveUserAvatars(u ...*User) *FileUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return fuo.RemoveUserAvatarIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (fuo *FileUpdateOne) Select(field string, fields ...string) *FileUpdateOne {
	fuo.fields = append([]string{field}, fields...)
	return fuo
}

// Save executes the query and returns the updated File entity.
func (fuo *FileUpdateOne) Save(ctx context.Context) (*File, error) {
	var (
		err  error
		node *File
	)
	fuo.defaults()
	if len(fuo.hooks) == 0 {
		node, err = fuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FileMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fuo.mutation = mutation
			node, err = fuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(fuo.hooks) - 1; i >= 0; i-- {
			if fuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = fuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (fuo *FileUpdateOne) SaveX(ctx context.Context) *File {
	node, err := fuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (fuo *FileUpdateOne) Exec(ctx context.Context) error {
	_, err := fuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fuo *FileUpdateOne) ExecX(ctx context.Context) {
	if err := fuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (fuo *FileUpdateOne) defaults() {
	if _, ok := fuo.mutation.UpdatedAt(); !ok {
		v := file.UpdateDefaultUpdatedAt()
		fuo.mutation.SetUpdatedAt(v)
	}
}

func (fuo *FileUpdateOne) sqlSave(ctx context.Context) (_node *File, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   file.Table,
			Columns: file.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: file.FieldID,
			},
		},
	}
	id, ok := fuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "File.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := fuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, file.FieldID)
		for _, f := range fields {
			if !file.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != file.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := fuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fuo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldUpdatedAt,
		})
	}
	if value, ok := fuo.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldDeletedAt,
		})
	}
	if fuo.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: file.FieldDeletedAt,
		})
	}
	if value, ok := fuo.mutation.Disk(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldDisk,
		})
	}
	if value, ok := fuo.mutation.Path(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldPath,
		})
	}
	if value, ok := fuo.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldType,
		})
	}
	if value, ok := fuo.mutation.Size(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: file.FieldSize,
		})
	}
	if value, ok := fuo.mutation.AddedSize(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: file.FieldSize,
		})
	}
	if fuo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   file.UserTable,
			Columns: []string{file.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   file.UserTable,
			Columns: []string{file.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fuo.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.PostsTable,
			Columns: []string{file.PostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: post.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.RemovedPostsIDs(); len(nodes) > 0 && !fuo.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.PostsTable,
			Columns: []string{file.PostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: post.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.PostsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.PostsTable,
			Columns: []string{file.PostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: post.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fuo.mutation.PagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.PagesTable,
			Columns: []string{file.PagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: page.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.RemovedPagesIDs(); len(nodes) > 0 && !fuo.mutation.PagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.PagesTable,
			Columns: []string{file.PagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: page.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.PagesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.PagesTable,
			Columns: []string{file.PagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: page.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fuo.mutation.UserAvatarsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.UserAvatarsTable,
			Columns: []string{file.UserAvatarsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.RemovedUserAvatarsIDs(); len(nodes) > 0 && !fuo.mutation.UserAvatarsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.UserAvatarsTable,
			Columns: []string{file.UserAvatarsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fuo.mutation.UserAvatarsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   file.UserAvatarsTable,
			Columns: []string{file.UserAvatarsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &File{config: fuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, fuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{file.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
