package service

import (
	"context"
	"testing/entity"
	"testing/repository"
)

type AuthorService struct {
	aRepo *repository.AuthorRepository
}

func NewAuthorService(aRepo *repository.AuthorRepository) *AuthorService {
	return &AuthorService{
		aRepo: aRepo,
	}
}

func (a *AuthorService) GetAllAuthor(ctx context.Context) ([]entity.Author, error) {
	return a.aRepo.GetAllAuthor(ctx)
}

func (a *AuthorService) AddAuthor(ctx context.Context, author entity.Author) (entity.Author, error) {
	err := a.aRepo.AddAuthor(ctx, author)
	if err != nil {
		return entity.Author{}, err
	}
	return author, nil
}

func (a *AuthorService) GetAuthorByID(ctx context.Context, id int) (entity.Author, error) {
	return a.aRepo.GetAuthorByID(ctx, id)
}

func (a *AuthorService) UpdateAuthor(ctx context.Context, author entity.Author) (entity.Author, error) {
	err := a.aRepo.UpdateAuthor(ctx, author)
	if err != nil {
		return entity.Author{}, err
	}

	return author, nil
}

func (a *AuthorService) DeleteAuthor(ctx context.Context, id int) error {
	return a.aRepo.DeleteAuthor(ctx, id)
}
