package usecase

import (
	"context"

	"github.com/pkg/errors"
	"local.com/go-clean-lambda/internal/domain"
	"local.com/go-clean-lambda/internal/logger"
)

var ErrInvalidInput error = errors.New("invalid input")

// DummyBo.
type DummyBo struct {
	ID   string
	Name string
	Attr string
}

// IsValid
//
//	@receiver cbo
//	@return bool
func (cbo *DummyBo) IsValid() bool {
	return cbo != nil && len(cbo.ID) > 0 && len(cbo.Name) > 0
}

// DummyUseCase.
type DummyUseCase struct {
	dummyRepo domain.DummyRepository
}

// NewDummyUseCase
//
//	@param dummyRepo
//	@return *DummyUseCase
func NewDummyUseCase(dummyRepo domain.DummyRepository) *DummyUseCase {
	return &DummyUseCase{
		dummyRepo: dummyRepo,
	}
}

// Get
//
//	@receiver uc
//	@param ctx
//	@param id
//	@return *DummyBo
//	@return error
func (uc *DummyUseCase) Get(ctx context.Context, id string) (*DummyBo, error) {
	if len(id) == 0 {
		return nil, errors.Wrap(ErrInvalidInput, "blank id")
	}
	entity, err := uc.dummyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "repository error with id: %s", id)
	}
	return uc.buildBo(entity), nil
}

// Add
//
//	@receiver uc
//	@param ctx
//	@param bo
//	@return *DummyBo
//	@return error
func (uc *DummyUseCase) Add(ctx context.Context, bo *DummyBo) (*DummyBo, error) {
	if bo == nil || !bo.IsValid() {
		return nil, errors.Wrapf(ErrInvalidInput, "invalid bo: %s", logger.Pretty(bo))
	}
	entity := uc.buildEntity(bo)
	entity, err := uc.dummyRepo.Insert(ctx, entity)
	if err != nil {
		return nil, errors.Wrapf(err, "repository error with bo: %s", logger.Pretty(bo))
	}
	return uc.buildBo(entity), nil
}

func (uc *DummyUseCase) Remove(ctx context.Context, id string) (*DummyBo, error) {
	if len(id) == 0 {
		return nil, errors.Wrap(ErrInvalidInput, "blank id")
	}
	entity, err := uc.dummyRepo.DeleteByID(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "repository error with id: %s", id)
	}
	return uc.buildBo(entity), nil
}

// buildEntity
//
//	@receiver uc
//	@param bo
//	@return *domain.Dummy
func (uc *DummyUseCase) buildEntity(bo *DummyBo) *domain.Dummy {
	if bo == nil {
		return nil
	}
	return &domain.Dummy{
		ID:       bo.ID,
		Name:     bo.Name,
		SomeAttr: bo.Attr,
	}
}

// buildBo
//
//	@receiver uc
//	@param entity
//	@return *DummyBo
func (uc *DummyUseCase) buildBo(entity *domain.Dummy) *DummyBo {
	if entity == nil {
		return nil
	}
	return &DummyBo{
		ID:   entity.ID,
		Name: entity.Name,
		Attr: entity.SomeAttr,
	}
}
