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

func (b *AuthorRepository) GetAllAuthor(ctx context.Context) ([]entity.Author, error) {
	var authorsTemp []entity.Author

	author, err := b.db.
		WithContext(ctx).
		Table("authors").
		Select("*").
		Where("deleted_at IS NULL").
		Rows()
	if err != nil {
		return []entity.Author{}, err
	}
	defer author.Close()

	for author.Next() {
		b.db.ScanRows(author, &authorsTemp)
	}

	// for _,value := range authorsTemp {
	// 	authors,err := b.db.
	// 	WithContext(ctx).
	// 	Table("authors").
	// 	Select("*").
	// 	Joins("inner join authors_authors on authors_authors.author_id = authors.id").
	// 	Where("authors_authors.id = ?",value.ID).
	// 	Scan(&authors)

	// 	value.Authors = append(value.Authors, authors)
	// }

	return authorsTemp, nil
}

func (b *AuthorRepository) AddAuthor(ctx context.Context, author entity.Author) error {
	err := b.db.
		WithContext(ctx).
		Create(&author).Error
	return err
}

func (b *AuthorRepository) GetAuthorByID(ctx context.Context, id int) (entity.Author, error) {
	var authorResult entity.Author

	err := b.db.
		WithContext(ctx).
		Table("authors").
		Where("id = ? AND deleted_at IS NULL", id).
		Find(&authorResult).Error
	if err != nil {
		return entity.Author{}, err
	}

	return authorResult, nil
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
