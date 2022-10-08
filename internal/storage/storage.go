package storage

import (
	"database/sql"
	"fmt"
	"golang-rest-api-server/internal/util"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func init() {
	DB = connectDb()
	applyDbMigrations(DB)
}

func connectDb() *sql.DB {
	fmt.Println("Establishing connection to database.")

	_db, openErr := sql.Open("sqlite", "./tasks.db")
	util.PanicError(openErr)

	pingError := _db.Ping()
	util.PanicError(pingError)

	fmt.Println("Database connection established.")

	return _db
}

func applyDbMigrations(db *sql.DB) {
	fmt.Println("Applying migrations.")

	taskTableInsertStmt := `
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT
		);
	`

	_, err := db.Exec(taskTableInsertStmt)
	util.PanicError(err)

	fmt.Println("Migration applied.")
}
