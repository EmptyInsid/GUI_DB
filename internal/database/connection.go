package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Init инициализирует подключение к БД
func (db *Database) Init(dsn string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	db.pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Can't connect with database: %v", err)
	}
	fmt.Println("Connected to database")
}

// Close закрывает пул соединений
func (db *Database) CloseDB() {
	if db.pool != nil {
		db.pool.Close()
	}
}

// Pool возвращает внутренний пул соединений (для использования напрямую, если нужно)
func (db *Database) Pool() *pgxpool.Pool {
	return db.pool
}

// GetEnvOrDefault возвращает значение переменной окружения или дефолтное значение
func GetEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
