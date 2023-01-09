package persistence

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/zhughes3/elliot/pkg/log"
)

// NewPostgresDB attempts to open a db handler to a postgres server
func NewPostgresDB(logger log.Logger, cfg dbConfig) (DB, error) {
	conn := newPostgresConnectionString(cfg)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return DB{}, fmt.Errorf("problem opening Postgres connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return DB{}, fmt.Errorf("problem verifying Postgres connection: %v", err)
	}

	logger.Infof("connected to postgres server at %s", cfg.Host)
	return NewDBFromConn(db), nil
}

func newPostgresConnectionString(cfg dbConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, "disable")
}
