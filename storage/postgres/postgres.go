package postgres

import (
	"errors"
	"fmt"
	"gopkg.in/pg.v5"
)

var schemaQueries []string = []string{
	`CREATE TABLE IF NOT EXISTS owners (username VARCHAR(32) NOT NULL UNIQUE PRIMARY KEY, password TEXT NOT NULL, salt TEXT NOT NULL, fullname VARCHAR(64) NOT NULL, email VARCHAR(64) NOT NULL, session_token VARCHAR(36), api_token VARCHAR(36), is_active BOOLEAN, is_admin BOOLEAN, signup_time TIMESTAMP, activation_time TIMESTAMP, last_login TIMESTAMP, last_login_host VARCHAR(39))`,
	`CREATE TABLE IF NOT EXISTS activations (token VARCHAR(36) NOT NULL UNIQUE PRIMARY KEY, username VARCHAR(32) NOT NULL REFERENCES owners, generation_time TIMESTAMP NOT NULL)`,
	`CREATE TABLE IF NOT EXISTS asnums (asnum INTEGER NOT NULL UNIQUE PRIMARY KEY, description VARCHAR(64) NOT NULL, username VARCHAR(32) NOT NULL REFERENCES owners)`,
	`CREATE TABLE IF NOT EXISTS prefixes (network CIDR NOT NULL UNIQUE PRIMARY KEY, description VARCHAR(64) NOT NULL, username VARCHAR(32) NOT NULL REFERENCES owners, parent CIDR REFERENCES prefixes)`,
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
			fmt.Println(err)
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
