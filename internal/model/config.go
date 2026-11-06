package model

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// Название приложения.
	NAME = "refresh-tokens-cleaner"

	// Имя переменной окружения определяющая режим работы приложения.
	MODE = "MODE"

	// Режим разработки.
	MODE_DEV = "dev"

	// Путь в директории конфигурации в производстве.
	CONFIG_DIRECTORY_PROD = "/etc"

	// Путь в директории конфигурации в разработке.
	CONFIG_DIRECTORY_DEV = "./config"

	// Имя файла конфигурации в производстве.
	CONFIG_FILENAME_PROD = "prod.yaml"

	// Имя файла конфигурации в разработке.
	CONFIG_FILENAME_DEV = "dev.yaml"
)

// Conf Глобальный экземпляр конфигурации.
var Conf *Config

// RedisConfig Конфигурация Redis.
type RedisConfig struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
	Db   int    `yaml:"db"`
}

// Config Конфигурация.
type Config struct {
	Redis           *RedisConfig `yaml:"redis"`
	TokenTTLMinutes int          `yaml:"token_time_to_live_minutes"`
}

// NewConfig Конфигурация.
func NewConfig() (*Config, error) {
	f, err := os.ReadFile(filepath.Join(configDir(), configFilename()))
	if err != nil {
		return nil, err
	}
	c := &Config{
		Redis: &RedisConfig{},
	}
	if err := yaml.Unmarshal(f, c); err != nil {
		return nil, err
	}
	return c, nil
}

// configDir Директоия конфигурации.
func configDir() string {
	m, o := os.LookupEnv(MODE)
	if o && m == MODE_DEV {
		return CONFIG_DIRECTORY_DEV
	}
	return filepath.Join(CONFIG_DIRECTORY_PROD, NAME)
}

// configFilename Имя файла конфигурации.
func configFilename() string {
	m, o := os.LookupEnv(MODE)
	if o && m == MODE_DEV {
		return CONFIG_FILENAME_DEV
	}
	return CONFIG_FILENAME_PROD
}
