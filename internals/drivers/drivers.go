package drivers

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB holds the database connection pool
// and provides methods to interact with the database.
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenConns = 10               // Maximum number of open connections to the database
const maxIdleConns = 5                // Maximum number of idle connections in the pool
const maxDbLifetime = 5 * time.Minute // Maximum lifetime of a connection in the pool (0 means no limit)

// Init initializes the database connection pool with the given driver and data source name (DSN).
func ConnectToSQL(dsn string) (*DB, error) {

	db, err := NewDatabase(dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(maxDbLifetime)

	dbConn.SQL = db

	// Test the connection to the database
	err = testConnection(db)
	if err != nil {
		return nil, err
	}
	return dbConn, nil

}

func testConnection(d *sql.DB) error {
	// Test the connection to the database
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
