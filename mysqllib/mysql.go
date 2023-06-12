package mysqllib

import (
	"database/sql"
	"fmt"
	"strings"
)

func TestDatabaseConnection(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully connected to database")
	}
}

func CreateDatabase(db *sql.DB, dbName string) {
	dbCreate := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)
	create, err := db.Query(dbCreate)
	if err != nil {
		panic(err.Error())
	}
	defer create.Close()
	_, err = db.Exec(fmt.Sprintf("USE %s", dbName))
	if err != nil {
		panic(err.Error())
	}
	logTableCreate := `CREATE TABLE IF NOT EXISTS log (
		id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
		log TEXT NOT NULL,
		created DATETIME DEFAULT NOW()
	)
`
	tableCreate, err := db.Query(logTableCreate)
	if err != nil {
		panic(err.Error())
	}
	defer tableCreate.Close()
}

func InsertLog(db *sql.DB, dbName string, log_lines []string) {
	_, err := db.Query(fmt.Sprintf("USE %s", dbName))
	if err != nil {
		panic(err.Error())
	}
	sql := "INSERT INTO log (log) VALUES "

	var values []interface{}
	for _, line := range log_lines {
		sql += "( ? ),"
		values = append(values, strings.TrimSpace(string(line)))
	}
	sql = sql[0 : len(sql)-1]
	insert, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}

	_, err = insert.Exec(values...)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}
