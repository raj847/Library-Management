package service

import (
	"context"
	"errors"
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

var (
	ErrDataMapping = errors.New("data already exist")
)

func (a *MappingService) GetAllMapping(ctx context.Context) ([]entity.Mapping, error) {
	return a.mRepo.GetAllMapping(ctx)
}

func (a *MappingService) AddMapping(ctx context.Context, mapping entity.Mapping) (entity.Mapping, error) {
	check, _ := a.mRepo.GetAllMapping(ctx)
	for _, v := range check {
		if v.AuthorReadID == mapping.AuthorReadID && v.BookReadID == mapping.BookReadID {
			return entity.Mapping{}, ErrDataMapping
		}
	}
	err := a.mRepo.AddMapping(ctx, mapping)
	if err != nil {
		return entity.Mapping{}, err
	}
	return mapping, nil
}
