package usecase_test

import (
	"context"
	nativeerr "errors"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"local.com/go-clean-lambda/internal/domain"
	"local.com/go-clean-lambda/internal/logger"
	"local.com/go-clean-lambda/internal/usecase"
)

const (
	invalidDummyID string = "invalid_dummy_id"
)

var errBadRepositoryAction error = nativeerr.New("mocked repo error")

func TestTraceableErrorLog(t *testing.T) {
	items := []*domain.Dummy{}
	repo := NewDummyMockRepository(items)
	u := usecase.NewDummyUseCase(repo)

	_, err := u.Get(context.TODO(), invalidDummyID)

	msg := "failed to test traceable error log"
	assertions := assert.New(t)
	assertions.NotNil(err, msg, "error not found")
	logger.Error("get with bad repo failed", err)
}

func TestDummyGetWithBadRepoReturnError(t *testing.T) {
	items := []*domain.Dummy{}
	repo := NewDummyMockRepository(items)
	u := usecase.NewDummyUseCase(repo)

	bo, err := u.Get(context.TODO(), invalidDummyID)

	msg := "get with bad repo didn't fail"
	assertions := assert.New(t)
	assertions.NotNil(err, msg, "error not found")
	assertions.Nil(bo, msg, "returned bo")
}

func TestDummyGetWithEmptyIDReturnError(t *testing.T) {
	items := []*domain.Dummy{}
	repo := NewDummyMockRepository(items)
	u := usecase.NewDummyUseCase(repo)

	bo, err := u.Get(context.TODO(), "")
	msg := "get with empty id didn't fail"
	assertions := assert.New(t)
	assertions.NotNil(err, msg, "error not found")
	assertions.True(errors.Is(err, usecase.ErrInvalidInput), msg, "error type")
	assertions.Nil(bo, msg, "returned bo")
}

func TestDummyGetWithIDReturnBo(t *testing.T) {
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

func NewDummyMockRepository(entities []*domain.Dummy) *DummyMockRepository {
	dmap := make(map[string]*domain.Dummy)
	for _, entity := range entities {
		dmap[entity.ID] = entity
	}
	return &DummyMockRepository{
		dmap: dmap,
	}
}

func (r *DummyMockRepository) GetByID(ctx context.Context, id string) (*domain.Dummy, error) {
	if id == "" {
		return nil, nil
	}
	if id == invalidDummyID {
		return nil, errors.Wrap(errBadRepositoryAction, "GetByID")
	}
	return r.dmap[id], nil
}

func (r *DummyMockRepository) Insert(ctx context.Context, dummy *domain.Dummy) (*domain.Dummy, error) {
	if dummy == nil || dummy.ID == "" {
		return nil, nil
	}
	if dummy.ID == invalidDummyID {
		return nil, errors.Wrap(errBadRepositoryAction, "Insert")
	}
	r.dmap[dummy.ID] = dummy
	return dummy, nil
}

func (r *DummyMockRepository) DeleteByID(ctx context.Context, id string) (*domain.Dummy, error) {
	if id == "" {
		return nil, nil
	}
	if id == invalidDummyID {
		return nil, errors.Wrap(errBadRepositoryAction, "DeleteByID")
	}
	dummy := r.dmap[id]
	delete(r.dmap, id)
	return dummy, nil
}
