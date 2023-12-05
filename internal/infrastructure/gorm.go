package infrastructure

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func EstablishConnection(DatabaseURI string) (*gorm.DB, error) {
	sqlDB, err := sql.Open("pgx", DatabaseURI)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, err
}
