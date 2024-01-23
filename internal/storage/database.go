package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type Database struct {
	Conn *pgx.Conn
}

var database *Database

func InitDatabase(connString string) {
	if database != nil {
		return
	}

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Println(err)
	}

	database = &Database{
		Conn: conn,
	}
}

func GetDatabase() *Database {
	return database
}

func (*Database) CloseConnection() {
	err := database.Conn.Close(context.Background())
	if err != nil {
		log.Error().Err(err).Msgf("")
	}
}
