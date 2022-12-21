package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"local.com/go-clean-lambda/internal/domain"
	"local.com/go-clean-lambda/internal/usecase"
)

func TestDummyGetWithIdReturnBo(t *testing.T) {
	item := &domain.Dummy{
		ID:       uuid.New().String(),
		Name:     "test_name",
		SomeAttr: "test_attr",
	}
	items := []*domain.Dummy{item}
	repo := NewDummyMockRepository(items)
	u := usecase.NewDummyUseCase(repo)

	bo, err := u.Get(context.TODO(), item.ID)
	msg := "failed to get dummy bo by id"
	assertions := assert.New(t)
	assertions.Nil(err, msg, "found error")
	assertions.NotNil(bo, msg, "nil bo")
	assertions.Equal(item.ID, bo.ID, msg, "id")
	assertions.Equal(item.Name, bo.Name, msg, "name")
	assertions.Equal(item.SomeAttr, bo.Attr, msg, "attr")
}

type DummyMockRepository struct {
	dmap map[string]*domain.Dummy
}

func NewDummyMockRepository(entities []*domain.Dummy) domain.DummyRepository {
	dmap := make(map[string]*domain.Dummy)
	for _, entity := range entities {
		dmap[entity.ID] = entity
	}
	return &DummyMockRepository{
		dmap: dmap,
	}
}

func (r *DummyMockRepository) GetById(ctx context.Context, id string) (*domain.Dummy, error) {
	if id == "" {
		return nil, nil
	}
	return r.dmap[id], nil
}

func (r *DummyMockRepository) Insert(ctx context.Context, dummy *domain.Dummy) (*domain.Dummy, error) {
	if dummy == nil || dummy.ID == "" {
		return nil, nil
	}
	r.dmap[dummy.ID] = dummy
	return dummy, nil
}

func (r *DummyMockRepository) DeleteById(ctx context.Context, id string) (*domain.Dummy, error) {
	if id == "" {
		return nil, nil
	}
	dummy := r.dmap[id]
	delete(r.dmap, id)
	return dummy, nil
}
