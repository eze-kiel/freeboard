package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zpatrick/go-config"
)

func Connect() (*sql.DB, error) {
	dbConf := config.NewYAMLFile("config/db-config.yaml")
	c := config.NewConfig([]config.Provider{dbConf})
	if err := c.Load(); err != nil {
		log.Fatal(err)
	}

	dbDriver := "mysql"
	dbUser, err := c.String("dbuser")
	if err != nil {
		log.Fatal(err)
	}

	dbPass, err := c.String("dbpass")
	if err != nil {
		log.Fatal(err)
	}

	dbName, err := c.String("dbname")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		return nil, err
	}
	return db, nil
}
