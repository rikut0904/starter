package config

import "os"

type Config struct {
	Port        string
	DatabaseDSN string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=app port=5432 sslmode=disable"
	}
	return Config{Port: port, DatabaseDSN: dsn}
}
