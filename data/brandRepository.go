package data

import (
	"database/sql"
	"errors"
)

type BrandRepository struct {
	db *sql.DB
}

var createBrand = `
INSERT INTO brand (name) VALUES (?);
`

func CreateBrandRepository(db *sql.DB) *BrandRepository {
	return &BrandRepository{
		db: db,
	}
}

func (r *BrandRepository) GetBrands() ([]string, error) {
	rows, err := r.db.Query("SELECT name FROM brand")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results = make([]string, 1)
	for rows.Next() {
		var name string
		rows.Scan(&name)
		results = append(results, name)
	}
	return results, nil
}

func (r *BrandRepository) CreateBrand(brand string) error {
	if len(brand) == 0 {
		return errors.New("brand must be a non empty string")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createBrand)

	if err != nil {
		return err
	}

	stmt.Exec(brand)
	stmt.Close()
	tx.Commit()

	return nil
}
