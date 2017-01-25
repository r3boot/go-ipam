package postgres

import (
	"errors"
	"gopkg.in/pg.v5"
)

var schemaQueries []string = []string{
	`CREATE TABLE IF NOT EXISTS owners (username TEXT NOT NULL UNIQUE PRIMARY KEY, fullname TEXT NOT NULL, email TEXT NOT NULL)`,
}

var db *pg.DB

type Config struct {
	Host     string
	User     string
	Pass     string
	Database string
}

func createSchema() error {
	var (
		q   string
		err error
	)
	for _, q = range schemaQueries {
		_, err = db.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}

func Connect(cfg interface{}) error {
	db = pg.Connect(&pg.Options{
		Addr:     cfg.(Config).Host,
		User:     cfg.(Config).User,
		Password: cfg.(Config).Pass,
		Database: cfg.(Config).Database,
	})

	if db == nil {
		return errors.New("Failed to connect to database")
	}

	createSchema()

	return nil
}
