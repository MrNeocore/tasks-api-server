package storage

import (
	"database/sql"
	"fmt"

	"github.com/MrNeocore/tasks-api-server/internal/util"

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

	tableExists := taskTableExists(db)

	if !tableExists {
		createTable(db)
	} else {
		applyMigrations(db)
	}

	fmt.Println("Migration applied.")
}

func taskTableExists(db *sql.DB) bool {
	stmt := `SELECT EXISTS (
		SELECT FROM 
			pg_tables
		WHERE 
			schemaname = 'public' AND 
			tablename  = 'tasks'
		);
	`

	row := db.QueryRow(stmt)

	var tableExists bool
	if err := row.Scan(&tableExists); err != nil {
		util.PanicError(err)
	}

	return tableExists
}

func createTable(db *sql.DB) {
	fmt.Println("Create `tasks` table.")

	taskTableInsertStmt := `
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			creationTime TIMESTAMP,
			shortTitle VARCHAR(32),
			title VARCHAR(256),
			description TEXT,
			tags TEXT[],
			category VARCHAR(64),
			priority SMALLINT,
			involvesOther BOOL,
			timeEstimate INTERVAL,
			dueDate TIMESTAMP,
			hardDeadline BOOL,
			reminder INTERVAL,
			repeats INTERVAL
		);
	`

	_, err := db.Exec(taskTableInsertStmt)
	util.PanicError(err)

	fmt.Println("`tasks` table created.")
}

func applyMigrations(db *sql.DB) {
	fmt.Println("Applying new schema to table `tasks`.")

	taskTableAlterStmt := `
		ALTER TABLE tasks 
		ADD COLUMN IF NOT EXISTS id TEXT PRIMARY KEY,
		ADD COLUMN IF NOT EXISTS creationTime TIMESTAMP,
		ADD COLUMN IF NOT EXISTS shortTitle VARCHAR(32),
		ADD COLUMN IF NOT EXISTS title VARCHAR(256),
		ADD COLUMN IF NOT EXISTS description TEXT,
		ADD COLUMN IF NOT EXISTS tags TEXT[],
		ADD COLUMN IF NOT EXISTS category VARCHAR(64),
		ADD COLUMN IF NOT EXISTS priority SMALLINT,
		ADD COLUMN IF NOT EXISTS involvesOther BOOL,
		ADD COLUMN IF NOT EXISTS timeEstimate INTERVAL,
		ADD COLUMN IF NOT EXISTS dueDate TIMESTAMP,
		ADD COLUMN IF NOT EXISTS hardDeadline BOOL,
		ADD COLUMN IF NOT EXISTS reminder INTERVAL,
		ADD COLUMN IF NOT EXISTS repeats INTERVAL
	`

	_, err := db.Exec(taskTableAlterStmt)
	util.PanicError(err)

	fmt.Println("New schema applied to table `tasks`.")
}
