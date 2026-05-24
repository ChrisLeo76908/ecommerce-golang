package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Conexion() *sql.DB {

	db, err := sql.Open(
		"mysql",
		"root:1234@tcp(localhost:3306)/ecommerce",
	)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
