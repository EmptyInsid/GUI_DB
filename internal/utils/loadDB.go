package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EmptyInsid/db_gui/internal/database"
)

func LoadDb(config *Config) (database.Service, error) {

	// Создаем экземпляр структуры Database
	db := &database.Database{}

	// Инициализация подключения
	db.Init(buildConnectionString(config))

	// Проверка подключения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Pool().Ping(ctx); err != nil {
		log.Fatalf("Error connection with bd: %v", err)
		return nil, err
	}

	return db, nil
}

func buildConnectionString(config *Config) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
		config.DBSSLMode,
	)
}
