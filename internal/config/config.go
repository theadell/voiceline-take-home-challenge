package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Port                   int
	Host                   string
	DatabaseURL            string
	MigrateDB              bool
	Oauth2ProviderName     string
	OAuth2DiscoveryUrl     string
	OAuth2ClientId         string
	OAuth2ClientSecret     string
	Oauth2CallbackUri      string
	OAuth2DeviceClientId   string
	OAuth2DeviceCientIdExt string
}

func LoadConfig() *Config {
	config := &Config{}

	config.Port = getEnvInt("PORT", 8080)
	config.Host = getEnvString("HOST", "localhost")

	config.MigrateDB = getEnvBool("MIGRATE_DB")
	config.DatabaseURL = getEnvString("DATABASE_URL", "./database.db")

	config.Oauth2ProviderName = getEnvString("OAUTH2_PROVIDER", "google")
	config.OAuth2DiscoveryUrl = getEnvString("OAUTH2_DISCOVERY_URL", "https://accounts.google.com/.well-known/openid-configuration")
	config.OAuth2ClientId = getEnvString("OAUTH2_CLIENT_ID", "44913867410-tdmgpcmgl9lm0sflp4vn8bfl36vbsf3v.apps.googleusercontent.com")
	config.OAuth2ClientSecret = getEnvString("OAUTH2_CLIENT_SECRET", "")

	config.Oauth2CallbackUri = getEnvString("OAUTH2_CALLBACK_URI", "http://localhost:8080/oauth2/google/callback")

	config.OAuth2DeviceClientId = getEnvString("OAUTH2_DEVICE_CLIENT_ID", "44913867410-2568om3gnua95hd47mrcn0tbao6iv6q4.apps.googleusercontent.com")
	config.OAuth2DeviceCientIdExt = getEnvString("OAUTH2_DEVICE_CLIENT_EXT", "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")

	flag.IntVar(&config.Port, "port", config.Port, "TCP Port to bind server to")
	flag.StringVar(&config.Host, "host", config.Host, "Network to bind to")

	flag.BoolVar(&config.MigrateDB, "migrate-db", config.MigrateDB, "Flag to enable DB migration on startup")
	flag.StringVar(&config.DatabaseURL, "database-url", config.DatabaseURL, "SQLITE Database URL")

	flag.StringVar(&config.Oauth2ProviderName, "oauth2-provider", config.Oauth2ProviderName, "Oauth2 provider name (e.g google)")
	flag.StringVar(&config.OAuth2DiscoveryUrl, "oauth2-discovery-url", config.OAuth2DiscoveryUrl, "Oauth2 provider discovery url")
	flag.StringVar(&config.OAuth2ClientId, "oauth2-client-id", config.OAuth2ClientId, "OAuth2 Client_id")
	flag.StringVar(&config.OAuth2ClientSecret, "oauth2-client-secret", config.OAuth2ClientSecret, "OAuth2 Client_secret")
	flag.StringVar(&config.Oauth2CallbackUri, "oauth2-redirect-uri", config.Oauth2CallbackUri, "OAuth2 Callback URI")

	flag.StringVar(&config.OAuth2DeviceClientId, "oauth2-device-client-id", config.OAuth2DeviceClientId, "OAuth2 device client id")
	flag.StringVar(&config.OAuth2DeviceCientIdExt, "oauth2-device-client-id-ext", config.OAuth2DeviceCientIdExt, "OAuth2 device client Id Extension")

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
