package repository

import (
	"context"
	"testing/entity"

	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db}
}

func (b *BookRepository) GetAllBook(ctx context.Context) ([]entity.BookRead, error) {
	var booksTemp []entity.Book

	book, err := b.db.
		WithContext(ctx).
		Table("books").
		Select("*").
		Where("deleted_at IS NULL").
		Rows()
	if err != nil {
		return []entity.BookRead{}, err
	}
	defer book.Close()

	for book.Next() {
		b.db.ScanRows(book, &booksTemp)
	}

	result := b.Relation(ctx, booksTemp)

	return result, nil

	// var books []entity.BookRead

	// err := b.db.WithContext(ctx).Table("books").Select("*").Preload("Authors").Where("deleted_at IS NULL").Find(&books).Error
	// if err != nil {
	// 	return nil, err
	// }

	// return books, nil
}

func (b *BookRepository) Relation(ctx context.Context, books []entity.Book) []entity.BookRead {
	result := []entity.BookRead{}
	for _, value := range books {
		mapping := []entity.Mapping{}
		b.db.WithContext(ctx).Select("*").Where("book_read_id =?", value.ID).Find(&mapping)

		authors := []entity.Author{}

		for _, mappings := range mapping {
			author := entity.Author{}
			b.db.First(&author, mappings.AuthorReadID)

			authors = append(authors, author)
		}

		bookRead := entity.BookRead{
			Model: gorm.Model{
				ID:        value.ID,
				CreatedAt: value.CreatedAt,
				DeletedAt: value.DeletedAt,
				UpdatedAt: value.UpdatedAt,
			},
			Title:         value.Title,
			ISBN:          value.ISBN,
			PublishedYear: value.PublishedYear,
			Authors:       authors,
		}

		result = append(result, bookRead)
	}
	return result
}

func (b *BookRepository) RelationID(ctx context.Context, books entity.Book) entity.BookRead {
	mapping := []entity.Mapping{}
	b.db.WithContext(ctx).Select("*").Where("book_read_id =?", books.ID).Find(&mapping)

	authors := []entity.Author{}
	for _, mappings := range mapping {
		author := entity.Author{}
		b.db.First(&author, mappings.AuthorReadID)

		authors = append(authors, author)
	}

	bookRead := entity.BookRead{
		Model: gorm.Model{
			ID:        books.ID,
			CreatedAt: books.CreatedAt,
			DeletedAt: books.DeletedAt,
			UpdatedAt: books.UpdatedAt,
		},
		Title:         books.Title,
		ISBN:          books.ISBN,
		PublishedYear: books.PublishedYear,
		Authors:       authors,
	}

	return bookRead
}

func (b *BookRepository) AddBook(ctx context.Context, book entity.Book) (entity.Book, error) {
	err := b.db.
		WithContext(ctx).
		Create(&book).Error
	if err != nil {
		return entity.Book{}, err
	}
	return book, nil
}

func (b *BookRepository) GetBookByID(ctx context.Context, id int) (entity.BookRead, error) {
	var bookResult entity.Book

	err := b.db.
		WithContext(ctx).
		Table("books").
		Where("id = ? AND deleted_at IS NULL", id).
		Find(&bookResult).Error
	if err != nil {
		return entity.BookRead{}, err
	}
	result := b.RelationID(ctx, bookResult)

	return result, nil
}

func (b *BookRepository) DeleteBook(ctx context.Context, id int) error {
	err := b.db.
		WithContext(ctx).
		Delete(&entity.Book{}, id).Error
	return err
}

func (b *BookRepository) UpdateBook(ctx context.Context, book entity.Book) error {
	err := b.db.
		WithContext(ctx).
		Table("books").
		Where("id = ?", book.ID).
		Updates(&book).Error
	return err
}
