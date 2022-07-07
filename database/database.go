package database

import (
	"database/sql"
	"fmt"
	"go-test/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Connection *sql.DB
}

var DatabaseInstance = &Database{}

func (d *Database) Load() bool {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.ConfigInstance.DbUser,
		config.ConfigInstance.DbPass,
		config.ConfigInstance.DbHost,
		config.ConfigInstance.DbPort,
		config.ConfigInstance.DbDatabase)

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatal(err)
		return false
	}

	d.Connection = db
	return true
}

func (d *Database) Version() string {
	var version string

	err := d.Connection.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	}

	return version
}
