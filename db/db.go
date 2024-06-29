package db

import (
	"database/sql"
	"log"
)

func ConfigureDB() *sql.DB {
	var (
		dbDriver = "mysql"
		dbSource = "root:root@tcp(mysql:3306)/images_db"
	)

	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("Can't configure db err=%s", err.Error())
		return nil
	}
	log.Println("successfully connected to db")
	return db
}
