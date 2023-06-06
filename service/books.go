package service

import (
	"context"
	"testing/entity"
	"testing/repository"
	"time"
)

type BookService struct {
	bRepo *repository.BookRepository
}

func NewBookService(bRepo *repository.BookRepository) *BookService {
	return &BookService{
		bRepo: bRepo,
	}
}

func (b *BookService) GetAllBook(ctx context.Context) ([]entity.BookRead, error) {
	return b.bRepo.GetAllBook(ctx)
}

func (b *BookService) AddBook(ctx context.Context, book entity.BookReq) (entity.Book, error) {
	t, _ := time.Parse("2006-01-02", book.PublishedYear)
	books, err := b.bRepo.AddBook(ctx, entity.Book{
		Title:         book.Title,
		PublishedYear: t,
		ISBN:          book.ISBN,
	})
	if err != nil {
		return entity.Book{}, err
	}
	return books, nil
}

func (b *BookService) GetBookByID(ctx context.Context, id int) (entity.BookRead, error) {
	return b.bRepo.GetBookByID(ctx, id)
}

func (b *BookService) UpdateBook(ctx context.Context, book entity.Book) (entity.Book, error) {
	err := b.bRepo.UpdateBook(ctx, book)
	if err != nil {
		return entity.Book{}, err
	}

	return book, nil
}

func (b *BookService) DeleteBook(ctx context.Context, id int) error {
	return b.bRepo.DeleteBook(ctx, id)
}
