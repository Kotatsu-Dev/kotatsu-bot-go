// ------------------------------------
// RR IT 2024
//
// ------------------------------------
// Базовый движок для Котацу бота

package config

const (
	CONFIG_URL_BASE = "URL_SITE"

	// ===========================================================
	// 							TELEGRAM CONFIG
	// ===========================================================

	//	Токен KotatsuBot - https://t.me/kotatsu_tg_bot
	CONFIG_BOT_TOKEN = "CONFIG_BOT_TOKEN"

	// Печать отладки в консоль красиво
	CONFIG_PRINT_LOG      = true
	CONFIG_PRINT_LOG_FILE = true

	// Сервер
	CONFIG_RELEASE_SERVER_PORT          = "8005"
	CONFIG_DEBUG_SERVERLESS_SERVER_PORT = "8006"

	// Отладка бота внутри
	IS_DEBUG_BOT = false

	// ID конфы техподдержки
	CONFIG_ID_CHAT_SUPPORT = -4017367566

	// ===========================================================
	// 							ROUTER CONFIG
	// ===========================================================

	// Режим отладки
	CONFIG_IS_DEBUG = false

	// Режим отладки Телеграмм Бота
	CONFIG_IS_BOT_DEBUG = false

	// Уровень отладки
	CONFIG_DEBUG_LEVEL = 1

	// Использование внутреннего сервера(для отладки)
	CONFIG_IS_DEBUG_SERVERLESS = true

	// ===========================================================
	// 							DB POSTGRE_SQL CONFIG
	// ===========================================================
	CONFIG_DB_HOST     = "localhost"
	CONFIG_DB_PORT     = "5432"
	CONFIG_DB_USER     = "CONFIG_DB_USER"
	CONFIG_DB_NAME     = "CONFIG_DB_NAME"
	CONFIG_DB_PASSWORD = "CONFIG_DB_PASSWORD"

	CONFIG_DB_IS_DEBUG = false
)
