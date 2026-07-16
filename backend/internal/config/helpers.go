package config

import (
	"database/sql"
	"fmt"

	"github.com/redis/go-redis/v9"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
	)
}

func (c *DatabaseConfig) OpenDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", c.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetConnMaxLifetime(c.ConnMaxLifetime)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *RedisConfig) NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.Addr(),
		Password: c.Password,
		DB:       c.DB,
		PoolSize: c.PoolSize,
	})
}

func (c *MinIOConfig) IsConfigured() bool {
	return c.Endpoint != "" && c.AccessKey != "" && c.SecretKey != ""
}

func (c *JWTConfig) IsConfigured() bool {
	return c.AccessSecret != "" && c.RefreshSecret != ""
}
