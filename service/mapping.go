package service

import (
	"context"
	"testing/entity"
	"testing/repository"
)

type MappingService struct {
	mRepo *repository.MappingRepository
}

func NewMappingService(mRepo *repository.MappingRepository) *MappingService {
	return &MappingService{
		mRepo: mRepo,
	}
}

func (a *MappingService) GetAllMapping(ctx context.Context) ([]entity.Mapping, error) {
	return a.mRepo.GetAllMapping(ctx)
}

func (a *MappingService) AddMapping(ctx context.Context, mapping entity.Mapping) (entity.Mapping, error) {
	err := a.mRepo.AddMapping(ctx, mapping)
	if err != nil {
		return entity.Mapping{}, err
	}
	return mapping, nil
}
