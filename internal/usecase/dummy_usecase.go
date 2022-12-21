package usecase

import (
	"context"
	"errors"
	"fmt"

	"local.com/go-clean-lambda/internal/domain"
	"local.com/go-clean-lambda/internal/logger"
)

// DummyBo
type DummyBo struct {
	ID   string
	Name string
	Attr string
}

// IsValid
//  @receiver cbo
//  @return bool
func (cbo *DummyBo) IsValid() bool {
	return cbo != nil && len(cbo.ID) > 0 && len(cbo.Name) > 0
}

// DummyUseCase
type DummyUseCase struct {
	dummyRepo domain.DummyRepository
}

// NewDummyUseCase
//  @param dummyRepo
//  @return *DummyUseCase
func NewDummyUseCase(dummyRepo domain.DummyRepository) *DummyUseCase {
	return &DummyUseCase{
		dummyRepo: dummyRepo,
	}
}

// Get
//  @receiver uc
//  @param ctx
//  @param id
//  @return *DummyBo
//  @return error
func (uc *DummyUseCase) Get(ctx context.Context, id string) (*DummyBo, error) {
	if len(id) == 0 {
		return nil, errors.New("failed to get entity by blank id")
	}
	entity, err := uc.dummyRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get entity from repo by id. %w", err)
	}
	if entity == nil {
		return nil, nil
	}
	bo, err := uc.buildBo(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("failed to get by id. id: %s. %w", id, err)
	}
	return bo, nil
}

// Add
//  @receiver uc
//  @param ctx
//  @param bo
//  @return *DummyBo
//  @return error
func (uc *DummyUseCase) Add(ctx context.Context, bo *DummyBo) (*DummyBo, error) {
	if bo == nil || !bo.IsValid() {
		return nil, fmt.Errorf("failed to add by bo. invalid bo: %s", logger.Pretty(bo))
	}
	entity, err := uc.buildEntity(ctx, bo)
	if err != nil {
		return nil, fmt.Errorf("failed to add by bo. bo: %s. %w", logger.Pretty(bo), err)
	}
	entity, err = uc.dummyRepo.Insert(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("failed to add by bo. bo: %s. %w", logger.Pretty(bo), err)
	}
	bo, err = uc.buildBo(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("failed to add by bo. bo: %s. %w", logger.Pretty(bo), err)
	}
	return bo, nil
}

func (uc *DummyUseCase) Remove(ctx context.Context, id string) (*DummyBo, error) {
	if len(id) == 0 {
		return nil, errors.New("failed to remove entity by blank id")
	}
	entity, err := uc.dummyRepo.DeleteById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to remove entity from repo by id. %w", err)
	}
	if entity == nil {
		return nil, nil
	}
	bo, err := uc.buildBo(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("failed to remove by id. id: %s. %w", id, err)
	}
	return bo, nil
}

// buildEntity
//  @receiver uc
//  @param ctx
//  @param bo
//  @return *domain.Dummy
//  @return error
func (uc *DummyUseCase) buildEntity(ctx context.Context, bo *DummyBo) (*domain.Dummy, error) {
	return &domain.Dummy{
		ID:       bo.ID,
		Name:     bo.Name,
		SomeAttr: bo.Attr,
	}, nil
}

// buildBo
//  @receiver uc
//  @param ctx
//  @param entity
//  @return *DummyBo
//  @return error
func (uc *DummyUseCase) buildBo(ctx context.Context, entity *domain.Dummy) (*DummyBo, error) {
	return &DummyBo{
		ID:   entity.ID,
		Name: entity.Name,
		Attr: entity.SomeAttr,
	}, nil
}
