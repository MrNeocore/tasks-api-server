package storage

import (
	"database/sql"
	"fmt"
	"golang-rest-api-server/internal/util"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func init() {
	DB = connectDb()
	applyDbMigrations(DB)
}

func connectDb() *sql.DB {
	fmt.Println("Establishing connection to database.")

	_db, openErr := sql.Open("pgx", "postgres://api:password@localhost:5555/api")
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
