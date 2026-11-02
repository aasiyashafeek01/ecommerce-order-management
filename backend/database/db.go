package database

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"time"

	mysql "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {

	config := mysql.Config{
		User:   getEnv("MYSQL_USER", "root"),
		Passwd: os.Getenv("MYSQL_PASSWORD"),
		Net:    "tcp",
		Addr: net.JoinHostPort(
			getEnv("MYSQL_HOST", "127.0.0.1"),
			getEnv("MYSQL_PORT", "3306"),
		),
		DBName:    getEnv("MYSQL_DATABASE", "ecommerce_db"),
		ParseTime: true,
		Loc:       time.Local,
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	return db, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
