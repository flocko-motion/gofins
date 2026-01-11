package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/flocko-motion/gofins/pkg/db/generated"
	"github.com/flocko-motion/gofins/pkg/files"
	"github.com/flocko-motion/gofins/pkg/log"
	_ "github.com/lib/pq"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var dbLogger *log.Logger

func init() {
	dbLogger = log.New("DB")
}

func logf(format string, args ...interface{}) {
	dbLogger.Printf(format+"\n", args...)
}

type DB struct {
	conn *sql.DB
}

var (
	globalDB *DB
	dbOnce   sync.Once
)

// ResetConnection resets the singleton so it can be retried.
// Only use this for retry logic during initialization.
func ResetConnection() {
	dbOnce = sync.Once{}
	globalDB = nil
}

// Db returns the global database connection singleton, initializing it on first call.
// This connection is shared across the entire application and should NEVER be closed.
// The connection is automatically managed and will be closed when the application exits.
// Returns nil if the connection cannot be established.
func Db() *DB {
	dbOnce.Do(func() {
		db, err := newDB()
		if err != nil {
			logf("Failed to initialize database connection: %v", err)
			return
		}
		globalDB = db
	})
	return globalDB
}

func genQ() *generated.Queries {
	return generated.New(Db().conn)
}

// newDB creates a new database connection
func newDB() (*DB, error) {
	// Try environment variables first (for Docker), then fall back to config file (for local dev)
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	user := getEnvOrDefault("DB_USER", "fins")
	dbname := getEnvOrDefault("DB_NAME", "fins")

	password := getEnvOrDefault("DB_PASSWORD", "")
	if password == "" {
		// Fallback to config file for local development
		var err error
		password, err = files.GetEnvValue("~/.fins/config/db.env", "POSTGRES_PASSWORD")
		if err != nil {
			return nil, fmt.Errorf("failed to read DB password from env or config file: %w", err)
		}
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)

	return &DB{conn: conn}, nil
}

// PrepareForShutdown closes the database connection during application shutdown.
// This should only be called once during graceful shutdown, not during normal operation.
// The singleton will remain closed after this call.
func PrepareForShutdown() error {
	db := Db()
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// Internal helper functions - not exported
func exec(query string, args ...interface{}) (sql.Result, error) {
	return Db().conn.Exec(query, args...)
}

func query(query string, args ...interface{}) (*sql.Rows, error) {
	return Db().conn.Query(query, args...)
}

// QueryRaw executes a raw SQL query and returns the rows
// This is exported for use by administrative commands
func QueryRaw(query string, args ...interface{}) (*sql.Rows, error) {
	return Db().conn.Query(query, args...)
}

// ExecRaw executes a raw SQL statement and returns the result
// This is exported for use by administrative commands
func ExecRaw(query string, args ...interface{}) (sql.Result, error) {
	return Db().conn.Exec(query, args...)
}

// ColumnInfo represents database column metadata
type ColumnInfo struct {
	TableName     string
	ColumnName    string
	DataType      string
	IsNullable    string
	ColumnDefault *string
}

// GetSchema returns the database schema information
func GetSchema() ([]ColumnInfo, error) {
	db := Db()
	query := `
		SELECT 
			c.table_name,
			c.column_name,
			c.data_type,
			c.is_nullable,
			c.column_default
		FROM information_schema.columns c
		WHERE c.table_schema = 'public'
		ORDER BY c.table_name, c.ordinal_position
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		if err := rows.Scan(&col.TableName, &col.ColumnName, &col.DataType, &col.IsNullable, &col.ColumnDefault); err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return columns, rows.Err()
}
