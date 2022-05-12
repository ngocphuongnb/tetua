package mockrepository

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/utils"
)

var FakeRepoErrors = map[string]error{}

type Repository[E entities.Entity] struct {
	Name     string
	entities []*E
	mu       sync.Mutex
}

type Filter struct {
	Search          string           `form:"search" json:"search"`
	Page            int              `form:"page" json:"page"`
	Limit           int              `form:"limit" json:"limit"`
	Sorts           []*entities.Sort `form:"orders" json:"orders"`
	IgnoreUrlParams []string         `form:"ignore_url_params" json:"ignore_url_params"`
	ExcludeIDs      []int            `form:"exclude_ids" json:"exclude_ids"`
}

func setEntityField[E entities.Entity](entity *E, field string, value interface{}) {
	reflect.ValueOf(entity).Elem().FieldByName(field).Set(reflect.ValueOf(value))
}

func idEQ[E entities.Entity](entity1, entity2 *E) bool {
	return getEntityField(entity1, "ID") == getEntityField(entity2, "ID")
}

func getEntityField[E entities.Entity](entity *E, field string) interface{} {
	r := reflect.ValueOf(entity)
	f := reflect.Indirect(r).FieldByName(field)
	return f.Interface()
}

func getEntityByField[E entities.Entity](name string, slice []*E, compareField string, compareValue interface{}) (*E, error) {
	foundEntities := utils.SliceFilter(slice, func(e *E) bool {
		return compareValue == getEntityField(e, compareField)
	})

	if len(foundEntities) == 0 {
		return nil, &entities.NotFoundError{Message: name + " not found with " + compareField + " = " + fmt.Sprintf("%v", compareValue)}
	}

	return foundEntities[0], nil
}

func ByID[E entities.Entity](ctx context.Context, name string, slice []*E, id int) (*E, error) {
	if ctx.Value("query_error") != nil {
		return nil, errors.New("ByID error")
	}

	return getEntityByField(name, slice, "ID", id)
}

func (m *Repository[E]) All(ctx context.Context) ([]*E, error) {
	return m.entities, nil
}

func (m *Repository[E]) ByID(ctx context.Context, id int) (*E, error) {
	return ByID(ctx, m.Name, m.entities, id)
}

func (m *Repository[E]) Create(ctx context.Context, entity *E) (*E, error) {
	if ctx.Value("create_error") != nil {
		return nil, errors.New("Error create " + m.Name)
	}

	if err, ok := FakeRepoErrors[m.Name+"_create"]; ok && err != nil {
		return nil, err
	}

	for _, e := range m.entities {
		if idEQ(e, entity) {
			return nil, errors.New(m.Name + " already exists")
		}
	}

	now := time.Now()
	m.mu.Lock()
	defer m.mu.Unlock()
	setEntityField(entity, "ID", len(m.entities)+1)
	setEntityField(entity, "CreatedAt", &now)
	setEntityField(entity, "UpdatedAt", &now)
	m.entities = append(m.entities, entity)

	return entity, nil
}

func (m *Repository[E]) Update(ctx context.Context, entity *E) (*E, error) {
	if ctx.Value("update_error") != nil {
		return nil, errors.New("Error save " + m.Name)
	}

	if err, ok := FakeRepoErrors[m.Name+"_update"]; ok && err != nil {
		return nil, err
	}

	found := false
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entities = utils.SliceMap(m.entities, func(e *E) *E {
		if idEQ(e, entity) {
			found = true
			return entity
		}
		return e
	})

	if !found {
		return nil, errors.New(m.Name + " not found")
	}

	return entity, nil
}

func (m *Repository[E]) DeleteByID(ctx context.Context, id int) error {
	if err, ok := FakeRepoErrors[m.Name+"_deleteByID"]; ok && err != nil {
		return err
	}

	found := false
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entities = utils.SliceFilter(m.entities, func(e *E) bool {
		if getEntityField(e, "ID") == id {
			found = true
		}
		return getEntityField(e, "ID") != id
	})

	if !found {
		return errors.New(m.Name + " not found")
	}

	return nil
}

func (m *Repository[E]) Find(ctx context.Context, filters ...*Filter) ([]*E, error) {
	if err, ok := FakeRepoErrors[m.Name+"_find"]; ok && err != nil {
		return nil, err
	}

	if len(filters) == 0 {
		return m.entities, nil
	}

	if filters[0].Page < 1 {
		filters[0].Page = 1
	}

	if filters[0].Limit < 1 {
		filters[0].Limit = 10
	}

	result := make([]*E, 0)
	filter := *filters[0]
	offset := (filter.Page - 1) * filter.Limit

	for index, e := range m.entities {
		if index < offset {
			continue
		}
		if index >= offset+filter.Limit {
			break
		}
		result = append(result, e)
	}

	return result, nil
}

func (m *Repository[E]) Paginate(ctx context.Context, filters ...*entities.Filter) (*entities.Paginate[E], error) {
	if err, ok := FakeRepoErrors[m.Name+"_paginate"]; ok && err != nil {
		return nil, err
	}

	result := make([]*E, 0)
	filter := *filters[0]
	offset := (filter.Page - 1) * filter.Limit

	for index, e := range m.entities {
		if index < offset {
			continue
		}
		if index >= offset+filter.Limit {
			break
		}
		result = append(result, e)
	}

	return &entities.Paginate[E]{
		PageCurrent: filter.Page,
		PageSize:    filter.Limit,
		Total:       len(m.entities),
		Data:        result,
	}, nil
}
