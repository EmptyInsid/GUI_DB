package main

import (
	"context"
	"fmt"
	"log"

	"github.com/EmptyInsid/db_gui/internal/db"
	"github.com/EmptyInsid/db_gui/internal/utils"
)

func main() {
	// Загрузка конфигурации
	config, err := utils.LoadConfig("../config/config.ini")
	if err != nil {
		log.Fatalf("Error connection with bd: %v", err)
	}

	// Инициализация базы данных
	db.InitDB(buildConnectionString(config))
	defer db.CloseDB()

	// Проверка подключения
	ctx := context.Background()
	if err := db.DB.Ping(ctx); err != nil {
		log.Fatalf("Error connection with bd: %v", err)
	}

	// инициализация приложения
	// создание главноего меню\меню входа
	// запуск приложения
}

func buildConnectionString(config *utils.Config) string {
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
