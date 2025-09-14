package database

import (
	"context"
	"time"

	"github.com/busnosh/go-utils/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) *pgxpool.Pool {
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		logger.Fatal("❌ Unable to parse DATABASE_URL: %v", err)
	}

	// Optional: configure pool limits
	poolConfig.MaxConns = 20
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	// Create connection pool
	dbpool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Fatal("❌ Unable to create connection pool: %v", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dbpool.Ping(ctx); err != nil {
		logger.Fatal("❌ Unable to ping database: %v", err)
	}

	logger.Info("✅ Connected to PostgreSQL successfully using URL!")
	return dbpool
}

// ClosePool closes the database pool
func ClosePool(dbpool *pgxpool.Pool) {
	if dbpool != nil {
		dbpool.Close()
	}
}
