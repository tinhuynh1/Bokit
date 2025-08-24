package database

import (
	"booking-svc/config"
	"booking-svc/pkg/logger"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PostgreSql *gorm.DB

func NewPostgres(cfg *config.PostgresConfig) (*gorm.DB, error) {
	dsn := config.GetPostgresDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open failed: %w", err)
	}

	// Set connection pool options
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping failed: %w", err)
	}

	PostgreSql, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open failed: %w", err)
	}
	return PostgreSql, nil
}

func Close() {
	if PostgreSql != nil {
		db, err := PostgreSql.DB()
		if err != nil {
			logger.L.Error("failed to get db", zap.Error(err))
			return
		}
		db.Close()
	}
}

func BeginTxn() *gorm.DB {
	return PostgreSql.Begin()
}

func RollbackTxn(tx *gorm.DB) {
	tx.Rollback()
}

func CommitTxn(tx *gorm.DB) error {
	return tx.Commit().Error
}
