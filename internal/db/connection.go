package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// InitDB инициализирует подключение к базе данных
func InitDB(dsn string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	DB, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Can't connect with database: %v", err)
	}
	fmt.Println("Connect with database")
}

// CloseDB закрывает пул подключений
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

// GetEnvOrDefault возвращает значение переменной окружения или дефолтное значение
func GetEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
