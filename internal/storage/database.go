package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
		log.Panic().Err(err).Msgf("connection to database refused")
	}

	database = &DatabaseStorage{
		Conn: conn,
	}

	err = database.startMigration(ctx)
	if err != nil {
		log.Panic().Err(err).Msgf("error start migration")
	}
}

func GetDatabaseStorage() *DatabaseStorage {
	return database
}

func (db *DatabaseStorage) startMigration(ctx context.Context) error {
	migration := `
			CREATE TABLE IF NOT EXISTS shortens (
				id				bigserial primary key, 
				original_url	varchar(255) not null unique, 
				short_url		varchar(255) not null,
			    user_id			uuid
		    )
		`

	_, err := db.Conn.Exec(ctx, migration)
	if err != nil {
		return err
	}

	return nil
}

func (db *DatabaseStorage) Close(ctx context.Context) error {
	db.Conn.Close()

	return nil
}

func (db *DatabaseStorage) AddShorten(ctx context.Context, shorten models.Shorten) error {
	args := pgx.NamedArgs{
		"originalURL": shorten.OriginalURL,
		"shortURL":    shorten.ShortURL,
		"userID":      sql.NullString{String: shorten.UserID, Valid: shorten.UserID != ""},
	}

	query := "INSERT INTO shortens (original_url, short_url, user_id) VALUES (@originalURL, @shortURL, @userID)"
	_, err := db.Conn.Exec(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return ErrUnique
		}
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
	query := "INSERT INTO shortens (original_url, short_url, user_id) VALUES (@originalURL, @shortURL, @userID)"
	for _, shorten := range shortens {
		args := pgx.NamedArgs{
			"originalURL": shorten.OriginalURL,
			"shortURL":    shorten.ShortURL,
			"userID":      sql.NullString{String: shorten.UserID, Valid: shorten.UserID != ""},
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

	return db.mapToShorten(row)
}

func (db *DatabaseStorage) GetByURL(ctx context.Context, URL string) (*models.Shorten, error) {
	sql := "SELECT id, original_url, short_url FROM shortens WHERE original_url = $1"
	row := db.Conn.QueryRow(ctx, sql, URL)

	return db.mapToShorten(row)
}

func (db *DatabaseStorage) GetUserURLS(ctx context.Context, userID string) ([]*models.Shorten, error) {
	sql := "SELECT id, original_url, short_url FROM shortens WHERE user_id = $1"
	rows, err := db.Conn.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	shortens, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*models.Shorten, error) {
		return db.mapToShorten(row)
	})
	if err != nil {
		return nil, err
	}

	return shortens, nil
}

func (db *DatabaseStorage) mapToShorten(row pgx.Row) (*models.Shorten, error) {
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
