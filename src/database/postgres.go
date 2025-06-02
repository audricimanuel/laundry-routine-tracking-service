package database

import (
	"context"
	"gin-boilerplate/src/tools"
	"github.com/jmoiron/sqlx"
)

const (
	DRIVER_POSTGRES = "postgres"
)

func InitializePostgresqlDatabaseSqlx(ctx context.Context, dsn string) *sqlx.DB {
	db := tools.NewSqlxDsn(ctx, DRIVER_POSTGRES, dsn)
	return db
}
