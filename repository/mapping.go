package repository

import (
	"context"
	"testing/entity"

	"gorm.io/gorm"
)

type MappingRepository struct {
	db *gorm.DB
}

func NewMappingRepository(db *gorm.DB) *MappingRepository {
	return &MappingRepository{db}
}

func (m *MappingRepository) GetAllMapping(ctx context.Context) ([]entity.Mapping, error) {
	var mappingsTemp []entity.Mapping

	mapping, err := m.db.
		WithContext(ctx).
		Table("mappings").
		Select("*").
		Rows()
	if err != nil {
		return []entity.Mapping{}, err
	}
	defer mapping.Close()

	for mapping.Next() {
		m.db.ScanRows(mapping, &mappingsTemp)
	}

	return mappingsTemp, nil
}

func (m *MappingRepository) AddMapping(ctx context.Context, mapping entity.Mapping) error {
	err := m.db.
		WithContext(ctx).
		Create(&mapping).Error
	return err
}
