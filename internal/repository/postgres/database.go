package postgres

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"

	"pet-shelter/internal/config"
)

func NewDB(cfg config.Config) (*sql.DB, error) {

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, err
	}

	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			break
		}

		log.Println("waiting for database...")
		time.Sleep(3 * time.Second)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
