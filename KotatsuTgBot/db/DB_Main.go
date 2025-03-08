package db

import (
	//Внутренние пакеты проекта
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/rr_debug"

	//Сторонние библиотеки
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	//Системные пакеты
	"fmt"
	"log"
	"os"
)

// ----------------------------------------------
//
// 				(Base) Общий функционал
//
// ----------------------------------------------

// Инициализация БД
func DB_Init() {
	db := DB_Database()

	// Миграция (настройка)
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Activity{})
	db.AutoMigrate(&AnimeRoulette{})
	db.AutoMigrate(&Request{})
}

// Функция коннекта к базе данных
func DB_Database() *gorm.DB {
	var logLevel logger.LogLevel

	if config.GetConfig().CONFIG_DB_IS_DEBUG {
		logLevel = logger.Error
	} else {
		logLevel = logger.Silent
	}

	// Установка уровня логирования в GORM
	errorLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logLevel,
		},
	)

	cfg := config.GetConfig()
	db_credentials := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", cfg.CONFIG_DB_HOST, cfg.CONFIG_DB_PORT, cfg.CONFIG_DB_USER, cfg.CONFIG_DB_NAME, cfg.CONFIG_DB_PASSWORD)

	db, err := gorm.Open(postgres.Open(db_credentials), &gorm.Config{
		Logger: errorLogger,
	})

	if err != nil {
		rr_debug.PrintLOG("DB_Main.go", "DB_Database", "gorm.Open", "Ошибка соединения с БД", err.Error())
	}
	return db
}
