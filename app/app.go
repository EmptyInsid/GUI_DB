package app

import (
	"log"

	"github.com/EmptyInsid/db_gui/internal/utils"
)

func Run() error {

	// Загрузка конфигурации
	config, err := utils.LoadConfig("../config/config.ini")
	if err != nil {
		log.Printf("Error connection with bd: %v", err)
		return err
	}

	//загрузка базы данных
	db, err := utils.LoadDb(config)
	defer db.CloseDB()
	if err != nil {
		log.Printf("Error connection with bd: %v", err)
		return err
	}

	utils.StartApp(db)

	return nil

}
