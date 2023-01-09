package persistence

import (
	"database/sql"
	"fmt"

	"github.com/knadh/koanf"
	"github.com/zhughes3/elliot/pkg/log"

	_ "github.com/lib/pq"
)

// DB is a simple wrapper around sql.DB
type DB struct {
	DB *sql.DB
}

type dbConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// NewDB opens a conn to database and verifies connection
func NewDB(logger log.Logger, cfg *koanf.Koanf) (DB, error) {
	dbCfg, err := parseDBConfiguration(logger, cfg)
	if err != nil {
		return DB{}, fmt.Errorf("problem parsing db configuration: %v", err)
	}

	return NewPostgresDB(logger, dbCfg)
}

// NewDBFromConn creates wrapper for sql.DB conn handler
func NewDBFromConn(db *sql.DB) DB {
	return DB{db}
}

// parseDBConfiguration attempts to create a dbConfig instance from koanf.Koanf configuration
func parseDBConfiguration(logger log.Logger, cfg *koanf.Koanf) (dbConfig, error) {
	host := cfg.String("DB_HOST")
	if len(host) == 0 {
		logger.Infof("no configuration found for database host, defaulting to: %s", "localhost")
		host = "localhost"
	}

	port := cfg.String("DB_PORT")
	if len(port) == 0 {
		logger.Infof("no configuration found for database port, defaulting to: %s", "5432")
		port = "5432"
	}

	user := cfg.String("DB_USER")
	if len(user) == 0 {
		return dbConfig{}, newRequiredEnvironmentVariableError("DB_USER")
	}

	password := cfg.String("DB_PASSWORD")
	if len(password) == 0 {
		return dbConfig{}, newRequiredEnvironmentVariableError("DB_PASSWORD")
	}

	name := cfg.String("DB_NAME")
	if len(name) == 0 {
		return dbConfig{}, newRequiredEnvironmentVariableError("DB_NAME")
	}

	return dbConfig{
		User:     user,
		Password: password,
		Name:     name,
		Host:     host,
		Port:     port,
	}, nil
}

func newRequiredEnvironmentVariableError(name string) error {
	return fmt.Errorf("missing required environment variable: %s", name)
}
