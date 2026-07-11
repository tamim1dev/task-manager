package database

import "github.com/jackc/pgx/v5/pgxpool"

type Database struct {
	Pool *pgxpool.Pool
}

var DB *Database = &Database{
	Pool: nil,
}
