package repositories

import "gorm.io/gorm"

type Repository struct {
	DB    *gorm.DB
	Debug bool
}

func NewRepository(db *gorm.DB, debug bool) *Repository {
	return &Repository{db, debug}
}
