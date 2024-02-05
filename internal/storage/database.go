package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pingachguk/ya-shortener/internal/models"
	"github.com/rs/zerolog/log"
)

type DatabaseStorage struct {
	Conn *pgxpool.Pool
}

var database *DatabaseStorage

func InitDatabase(ctx context.Context, connString string) {
	if database != nil {
		return
	}

	conn, err := pgxpool.New(context.Background(), connString)
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
			id				bigserial primary key, 
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
	db.Conn.Close()

	return nil
}

func (db *DatabaseStorage) AddShorten(ctx context.Context, shorten models.Shorten) error {
	query := "INSERT INTO shortens (original_url, short_url) VALUES ($1, $2)"
	_, err := db.Conn.Exec(ctx, query, shorten.OriginalURL, shorten.ShortURL)
	if err != nil {
		return err
	}

	return nil
}

func (db *DatabaseStorage) AddBatchShorten(ctx context.Context, shortens []models.Shorten) error {
	tx, err := db.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	query := "INSERT INTO shortens (original_url, short_url) VALUES (@originalURL, @shortURL)"
	for _, shorten := range shortens {
		args := pgx.NamedArgs{
			"originalURL": shorten.OriginalURL,
			"shortURL":    shorten.ShortURL,
		}
		batch.Queue(query, args)
	}

	results := tx.SendBatch(ctx, batch)

	for range shortens {
		_, err := results.Exec()
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	if err := results.Close(); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (db *DatabaseStorage) GetByShort(ctx context.Context, short string) (*models.Shorten, error) {
	sql := "SELECT id, original_url, short_url FROM shortens WHERE short_url = $1"
	row := db.Conn.QueryRow(ctx, sql, short)
	shorten := &models.Shorten{}

	err := row.Scan(&shorten.UUID, &shorten.OriginalURL, &shorten.ShortURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return shorten, nil
}
