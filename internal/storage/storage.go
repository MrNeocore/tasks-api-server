package storage

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/MrNeocore/tasks-api-server/internal/util"
	"github.com/MrNeocore/tasks-api-server/task"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var HOST = util.GetOrElse(os.LookupEnv, "DB_HOST", "localhost")
var PORT = util.GetOrElse(os.LookupEnv, "DB_PORT", "5432")
var USER = util.GetOrElse(os.LookupEnv, "DB_USER", "api")
var PWD = os.Getenv("DB_PWD")
var NAME = util.GetOrElse(os.LookupEnv, "DB_NAME", "api")

var DB *sql.DB

var minVersionStr = ">=12"
var minVersion, _ = semver.NewConstraint(minVersionStr)

func init() {
	DB = connectDb()
	checkVersionAtLeast(DB, *minVersion)
	applyDbMigrations(DB)
}

func connectDb() *sql.DB {
	fmt.Println("Establishing connection to database.")

	connectionString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", USER, PWD, HOST, PORT, NAME)
	_db, openErr := sql.Open("pgx", connectionString)
	util.PanicError(openErr)

	pingError := _db.Ping()
	util.PanicError(pingError)

	fmt.Println("Database connection established.")

	return _db
}

func checkVersionAtLeast(db *sql.DB, minVersion semver.Constraints) {
	var currentVersion string
	err := db.QueryRow("SELECT split_part(version(), ' ', 2)").Scan(&currentVersion)
	if err != nil {
		util.PanicError(err)
	}

	if !minVersion.Check(semver.MustParse(currentVersion)) {
		err = fmt.Errorf("postgres version is too old: %v. Requires: %v", currentVersion, minVersionStr)
		util.PanicError(err)
	}
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

	columns := getTaskColumnNameAndTypes()
	createStmts := make([]string, len(columns))

	for i, column := range columns {
		alterStmt := fmt.Sprintf("%v %v", column.n, column.t)
		createStmts[i] = alterStmt
	}

	taskTableInsertStmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS tasks (%v);", strings.Join(createStmts, ","))

	_, err := db.Exec(taskTableInsertStmt)
	util.PanicError(err)

	fmt.Println("`tasks` table created.")
}

func applyMigrations(db *sql.DB) {
	fmt.Println("Applying new schema to table `tasks`.")

	columns := getTaskColumnNameAndTypes()
	alterStmts := make([]string, len(columns))

	for i, column := range columns {
		alterStmt := fmt.Sprintf("ADD COLUMN IF NOT EXISTS %v %v", column.n, column.t)
		alterStmts[i] = alterStmt
	}

	taskTableAlterStmt := fmt.Sprintf("ALTER TABLE tasks %v", strings.Join(alterStmts, ","))

	_, err := db.Exec(taskTableAlterStmt)
	util.PanicError(err)

	fmt.Println("New schema applied to table `tasks`.")
}

type PgColumn struct {
	n string
	t string
}

// Move to task ?
func getTaskColumnNameAndTypes() []PgColumn {
	taskType := reflect.TypeOf(&task.Task{}).Elem()

	columns := make([]PgColumn, taskType.NumField())

	for i := 0; i < len(columns); i++ {
		f := taskType.Field(i)
		fieldName := string(f.Tag.Get("json"))
		fieldType := string(f.Tag.Get("pgtype"))
		columns[i] = PgColumn{fieldName, fieldType}
	}

	return columns
}
