package repository

import (
	"github.com/epysqyli/anchors-backend/domain"
	"gorm.io/gorm"
)

type TagRepository struct {
	database *gorm.DB
}

func NewTagRepository(db *gorm.DB) domain.TagRepository {
	return &TagRepository{
		database: db,
	}
}

func (tr *TagRepository) Create(tag *domain.Tag) (domain.Tag, error) {
	res := tr.database.Create(tag)
	if res.Error != nil {
		return *tag, res.Error
	}

	return *tag, nil
}

func (tr *TagRepository) FetchAll() []domain.Tag {
	tags := []domain.Tag{}
	tr.database.Find(&tags)

	return tags
}

func (tr *TagRepository) FetchById(ID string) domain.Tag {
	tag := domain.Tag{}
	tr.database.First(&tag, ID)

	return tag
}

func (tr *TagRepository) FetchByName(name string) domain.Tag {
	tag := domain.Tag{}
	tr.database.Where(&domain.Tag{Name: name}).First(&tag)

	return tag
}

func (tr *TagRepository) DeleteByID(ID string) error {
	tx := tr.database.Delete(&domain.Tag{}, ID)
	return tx.Error
}
