package presentation

import (
	"flag"
	"os"
)

var (
	ServerEndpoint   = flag.String("a", "localhost:8080", "Server endpoint")
	DatabaseURI      = flag.String("d", "postgresql://admin:admin@localhost:5432/postgres", "Database URI")
	AccrualSystemURL = flag.String("r", "localhost:8080", "Accrual system address")
	SlugSize         = 8
	TokenExp         = 3600
	SecretKey        = "supersecretkey"
)

func ParseFlags() {
	flag.Parse()

	if envServerAddr, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		*ServerEndpoint = envServerAddr
	}

	if envDatabaseURI, ok := os.LookupEnv("DATABASE_URI"); ok {
		*DatabaseURI = envDatabaseURI
	}

	if envAccrualSystemURL, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); ok {
		*AccrualSystemURL = envAccrualSystemURL
	}
}
