package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var migrations = [][]string{
	{
		`CREATE TABLE migration (
			version integer
		)`,
		`INSERT INTO migration values (0)`,
		`CREATE TABLE brand (
			id integer primary key,
			name text not null
		)`,
		`CREATE TABLE shortened (
			id integer primary key,
			name text not null, 
			brand_id integer not null,
			FOREIGN KEY(brand_id) REFERENCES brand(id)
		)`,
		`CREATE INDEX brand_and_name ON shortened (brand_id, name)`,
	},
}

func Migrate() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:./foo.db?_foreign_keys=true")
	if err != nil {
		return nil, err
	}

	indexStmt := `
		SELECT name from sqlite_master
		WHERE type='table'
		and name='migration'
	`

	res := db.QueryRow(indexStmt)
	var migration string
	err = res.Scan(&migration)

	if err != nil && err == sql.ErrNoRows {
		log.Printf("INFO: Database has not yet been initialised")
		err = migrate(0, db)
		if err != nil {
			log.Fatal("Unable to migrate database ", err)
		}
	} else if err != nil {
		log.Fatal("Unable to migrate database ", err)
	} else {
		res = db.QueryRow(`SELECT version from migration`)
		var version int
		err = res.Scan(&version)
		if err != nil {
			log.Fatal("Unable to migrate database ", err)
		}
		err = migrate(version, db)
		if err != nil {
			log.Fatal("Unable to migrate database ", err)
		}
	}

	return db, nil
}

func migrate(version int, db *sql.DB) error {
	for index, migrationSet := range migrations {
		if index < version {
			continue
		}
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		for _, individual := range migrationSet {
			_, err := tx.Exec(individual)

			if err != nil {
				log.Printf("Rolling back as the migration (%d) has failed", index)
				rollbackErr := tx.Rollback()
				if rollbackErr != nil {
					log.Fatal(err)
				}
				return err
			}
		}
		// set the transaction level
		stmt, _ := tx.Prepare("UPDATE migration SET version = ?")
		stmt.Exec(index + 1)
		stmt.Close()
		tx.Commit()
	}
	return nil
}
