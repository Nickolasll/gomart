package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerEndpoint   string
	DatabaseURI      string
	AccrualSystemURL string
	TokenExp         int
	SecretKey        string
}

var (
	serverEndpoint   = flag.String("a", "localhost:8080", "Server endpoint")
	databaseURI      = flag.String("d", "postgresql://admin:admin@localhost:5432/postgres", "Database URI")
	accrualSystemURL = flag.String("r", "localhost:8080", "Accrual system address")
)

func GetConfig() Config {
	flag.Parse()

	if envServerAddr, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		*serverEndpoint = envServerAddr
	}

	if envDatabaseURI, ok := os.LookupEnv("DATABASE_URI"); ok {
		*databaseURI = envDatabaseURI
	}

	if envAccrualSystemURL, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); ok {
		*accrualSystemURL = envAccrualSystemURL
	}
	return Config{
		ServerEndpoint:   *serverEndpoint,
		DatabaseURI:      *databaseURI,
		AccrualSystemURL: *accrualSystemURL,
		TokenExp:         3600,
		SecretKey:        "supersecretkey",
	}
}
