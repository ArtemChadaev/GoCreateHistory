package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	// Основные настройки приложения
	Port       string `mapstructure:"PORT"`
	ComfyUIUrl string `mapstructure:"COMFU_UI_URL"`

	// Настройки Postgres
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	// Настройки Redis
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
}

func Load() (*Config, error) {
	v := viper.New()

	// 1. Разрешаем Viper автоматически искать переменные в системе
	v.AutomaticEnv()

	// 2. Список ключей, которые Docker Compose прокидывает в контейнер
	// Мы явно "привязываем" их, чтобы Unmarshal смог наполнить структуру
	keys := []string{
		"PORT",
		"COMFU_UI_URL",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_NAME", "DB_PASSWORD",
		"REDIS_HOST", "REDIS_PORT", "REDIS_PASSWORD",
	}

	for _, key := range keys {
		if err := v.BindEnv(key); err != nil {
			return nil, fmt.Errorf("ошибка привязки переменной %s: %w", key, err)
		}
	}

	var cfg Config
	// 3. Заполняем структуру данными из привязанных переменных окружения
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("ошибка распаковки конфига: %w", err)
	}

	return &cfg, nil
}
