// ------------------------------------
// RR IT 2024
//
// ------------------------------------
// Базовый движок для Котацу бота

package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	CONFIG_URL_BASE string

	CONFIG_BOT_TOKEN                    string
	CONFIG_PRINT_LOG                    bool
	CONFIG_PRINT_LOG_FILE               bool
	CONFIG_RELEASE_SERVER_PORT          string
	CONFIG_DEBUG_SERVERLESS_SERVER_PORT string
	IS_DEBUG_BOT                        bool
	CONFIG_ID_CHAT_SUPPORT              int64

	CONFIG_IS_DEBUG            bool
	CONFIG_IS_BOT_DEBUG        bool
	CONFIG_DEBUG_LEVEL         int
	CONFIG_IS_DEBUG_SERVERLESS bool

	CONFIG_DB_HOST     string
	CONFIG_DB_PORT     string
	CONFIG_DB_USER     string
	CONFIG_DB_NAME     string
	CONFIG_DB_PASSWORD string
	CONFIG_DB_IS_DEBUG bool
}

var config *Config = nil

func GetConfig() Config {
	if config != nil {
		return *config
	}

	data, err := os.ReadFile("./config.toml")
	if err != nil {
		panic(err.Error())
	}

	var tmpConfig Config
	err = toml.Unmarshal(data, &tmpConfig)
	if err != nil {
		panic(err.Error())
	}

	config = &tmpConfig
	return tmpConfig
}
