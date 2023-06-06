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
	// var booksTemp []entity.BookRead

	// book, err := b.db.
	// 	WithContext(ctx).
	// 	Table("books").
	// 	Select("*").
	// 	Where("deleted_at IS NULL").
	// 	Rows()
	// if err != nil {
	// 	return []entity.BookRead{}, err
	// }
	// defer book.Close()

	// for book.Next() {
	// 	b.db.ScanRows(book, &booksTemp)
	// }

	// for _, value := range booksTemp {
	// 	authors := entity.Author{}
	// 	b.db.
	// 		WithContext(ctx).
	// 		Table("authors").
	// 		Select("*").
	// 		Joins("inner join mappings on mappings.authors_id = authors.id").
	// 		Where("mappings.books_id = ?", value.ID).
	// 		Scan(&authors)

	// 	value.Authors = append(value.Authors, authors)
	// }

	// return booksTemp, nil

	var books []entity.BookRead

	err := b.db.WithContext(ctx).Table("books").Select("*").Preload("Authors").Where("deleted_at IS NULL").Find(&books).Error
	if err != nil {
		return nil, err
	}

	return books, nil
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
	var bookResult entity.BookRead

	err := b.db.
		WithContext(ctx).
		Table("books").
		Where("id = ? AND deleted_at IS NULL", id).
		Find(&bookResult).Error
	if err != nil {
		return entity.BookRead{}, err
	}
	authors := entity.Author{}
	b.db.
		WithContext(ctx).
		Table("authors").
		Select("*").
		Joins("inner join mappings on mappings.authors_id = authors.id").
		Where("mappings.books_id = ?", bookResult.ID).
		Scan(&authors)

	bookResult.Authors = append(bookResult.Authors, &authors)

	return bookResult, nil
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
