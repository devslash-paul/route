package data

import "database/sql"

type RedirectRepository struct {
	db *sql.DB
}

func Create(db *sql.DB) *RedirectRepository {
	return &RedirectRepository{
		db: db,
	}
}

func (*RedirectRepository) createRow() {

}
