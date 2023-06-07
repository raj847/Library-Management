package repository

import (
	"context"
	"testing/entity"

	"gorm.io/gorm"
)

type AuthorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{db}
}

func (b *AuthorRepository) GetAllAuthor(ctx context.Context) ([]entity.AuthorRead, error) {
	var authorsTemp []entity.Author

	author, err := b.db.
		WithContext(ctx).
		Table("authors").
		Select("*").
		Where("deleted_at IS NULL").
		Rows()
	if err != nil {
		return []entity.AuthorRead{}, err
	}
	defer author.Close()

	for author.Next() {
		b.db.ScanRows(author, &authorsTemp)
	}

	result := b.Relation(ctx, authorsTemp)

	return result, nil
}

func (b *AuthorRepository) AddAuthor(ctx context.Context, author entity.Author) error {
	err := b.db.
		WithContext(ctx).
		Create(&author).Error
	return err
}

func (b *AuthorRepository) GetAuthorByID(ctx context.Context, id int) (entity.AuthorRead, error) {
	var authorResult entity.Author

	err := b.db.
		WithContext(ctx).
		Table("authors").
		Where("id = ? AND deleted_at IS NULL", id).
		Find(&authorResult).Error
	if err != nil {
		return entity.AuthorRead{}, err
	}

	result := b.RelationID(ctx, authorResult)

	return result, nil
}

func (b *AuthorRepository) DeleteAuthor(ctx context.Context, id int) error {
	err := b.db.
		WithContext(ctx).
		Delete(&entity.Author{}, id).Error
	return err
}

func (b *AuthorRepository) UpdateAuthor(ctx context.Context, author entity.Author) error {
	err := b.db.
		WithContext(ctx).
		Table("authors").
		Where("id = ?", author.ID).
		Updates(&author).Error
	return err
}

func (b *AuthorRepository) Relation(ctx context.Context, authors []entity.Author) []entity.AuthorRead {
	result := []entity.AuthorRead{}
	for _, value := range authors {
		mapping := []entity.Mapping{}
		b.db.WithContext(ctx).Select("*").Where("author_read_id =?", value.ID).Find(&mapping)

		books := []entity.Book{}

		for _, mappings := range mapping {
			book := entity.Book{}
			b.db.First(&book, mappings.BookReadID)

			books = append(books, book)
		}

		authorRead := entity.AuthorRead{
			Model: gorm.Model{
				ID:        value.ID,
				CreatedAt: value.CreatedAt,
				DeletedAt: value.DeletedAt,
				UpdatedAt: value.UpdatedAt,
			},
			Name:    value.Name,
			Country: value.Country,
			Books:   books,
		}

		result = append(result, authorRead)
	}
	return result
}

func (b *AuthorRepository) RelationID(ctx context.Context, authors entity.Author) entity.AuthorRead {
	mapping := []entity.Mapping{}
	b.db.WithContext(ctx).Select("*").Where("author_read_id =?", authors.ID).Find(&mapping)

	books := []entity.Book{}
	for _, mappings := range mapping {
		book := entity.Book{}
		b.db.First(&book, mappings.BookReadID)

		books = append(books, book)
	}

	authorRead := entity.AuthorRead{
		Model: gorm.Model{
			ID:        authors.ID,
			CreatedAt: authors.CreatedAt,
			DeletedAt: authors.DeletedAt,
			UpdatedAt: authors.UpdatedAt,
		},
		Name:    authors.Name,
		Country: authors.Country,
		Books:   books,
	}

	return authorRead
}
