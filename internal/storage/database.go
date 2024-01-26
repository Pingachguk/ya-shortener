package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pingachguk/ya-shortener/internal/models"
	"github.com/rs/zerolog/log"
)

type DatabaseStorage struct {
	Conn *pgx.Conn
}

var database *DatabaseStorage

func InitDatabase(ctx context.Context, connString string) {
	if database != nil {
		return
	}

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Println(err)
	}

	database = &DatabaseStorage{
		Conn: conn,
	}

	err = database.startMigrations(ctx)
	if err != nil {
		log.Panic().Err(err).Msgf("")
	}
}

func GetDatabaseStorage() *DatabaseStorage {
	return database
}

func (db *DatabaseStorage) startMigrations(ctx context.Context) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS shortens (
			id				serial primary key, 
			original_url	varchar(255) not null, 
			short_url		varchar(255) not null
    	)`,
	}

	for _, query := range queries {
		_, err := db.Conn.Exec(ctx, query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DatabaseStorage) Close(ctx context.Context) error {
	return db.Conn.Close(ctx)
}

func (db *DatabaseStorage) AddShorten(ctx context.Context, shorten models.Shorten) error {
	sql := "INSERT INTO shortens (original_url, short_url) VALUES ($1, $2)"
	_, err := db.Conn.Exec(ctx, sql, shorten.OriginalURL, shorten.ShortURL)
	if err != nil {
		return err
	}

	return nil
}

func (db *DatabaseStorage) GetByShort(ctx context.Context, short string) (*models.Shorten, error) {
	sql := "SELECT id, original_url, short_url FROM shortens WHERE short = $1"
	row := db.Conn.QueryRow(ctx, sql, short)
	shorten := &models.Shorten{}

	err := row.Scan(shorten.UUID, shorten.OriginalURL, shorten.ShortURL)
	if err != nil {
		return nil, err
	}
	return shorten, nil
}
