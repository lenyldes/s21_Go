package shared

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// postgres
	PgUser     string
	PgPassword string
	PgDatabase string
	PgHost     string
	PgPort     string
	// app
	AppPort string
	// jwt tokens
	AccessTokenKey  string
	RefreshTokenKey string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Leeway          time.Duration
	Issuer          string
}

func NewConfig() *Config {
	godotenv.Load(".env")
	return &Config{
		PgUser:     getEnv("POSTGRES_USER"),
		PgPassword: getEnv("POSTGRES_PASSWORD"),
		PgDatabase: getEnv("POSTGRES_DB"),
		PgHost:     getEnv("POSTGRES_HOST"),
		PgPort:     getEnv("POSTGRES_PORT"),

		AppPort: getEnv("APP_PORT"),

		AccessTokenKey:  getEnv("ACCESS_TOKEN_KEY"),
		RefreshTokenKey: getEnv("REFRESH_TOKEN_KEY"),
		AccessTokenTTL:  getDuration("ACCESS_TOKEN_TTL_MINS", time.Minute),
		RefreshTokenTTL: getDuration("REFRESH_TOKEN_TTL_DAYS", time.Hour*24),
		Leeway:          getDuration("LEEWAY_MINS", time.Minute),
		Issuer:          getEnv("ISSUER"),
	}
}

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("env value for key <%s> not set", key))
	}
	return value
}

func getDuration(key string, unit time.Duration) time.Duration {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("env value for key <%s> not set", key))
	}
	number, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("value for key %s not a number", key))
	}
	return time.Duration(number) * unit
}
