package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Port        int
	Host        string
	DatabaseURL string
	MigrateDB   bool
}

func LoadConfig() *Config {
	config := &Config{}

	config.Port = getEnvInt("PORT", 8080)
	config.Host = getEnvString("HOST", "localhost")

	config.MigrateDB = getEnvBool("MIGRATE_DB")
	config.DatabaseURL = getEnvString("DATABASE_URL", "./database.db")

	flag.IntVar(&config.Port, "port", config.Port, "TCP Port to bind server to")
	flag.StringVar(&config.Host, "host", config.Host, "Network to bind to")

	flag.BoolVar(&config.MigrateDB, "migrate-db", config.MigrateDB, "Flag to enable DB migration on startup")
	flag.StringVar(&config.DatabaseURL, "database-url", config.DatabaseURL, "SQLITE Database URL")

	flag.Parse()

	return config
}

func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
func getEnvBool(key string) bool {
	vString := getEnvString(key, "false")
	val, err := strconv.ParseBool(vString)
	if err != nil {
		return false
	}
	return val
}
