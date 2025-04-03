package repository

import (
	"database/sql"
)

type PicRepository struct {
	db *sql.DB
}

func NewPicRepository(db *sql.DB) *PicRepository {
	return &PicRepository{
		db: db,
	}
}
