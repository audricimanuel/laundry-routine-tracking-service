package database

import (
	"context"
	"fmt"
	"gin-boilerplate/src/config"
	"github.com/jmoiron/sqlx"
)

type (
	DBCollection struct {
		PostgresDBSqlx *sqlx.DB
	}
)

func NewDatabaseCollection(cfg config.Config) DBCollection {
	ctx := context.Background()

	// postgres
	postgresDBConfig := cfg.DataSource.PostgresDBConfig
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		postgresDBConfig.Host, postgresDBConfig.User, postgresDBConfig.Password, postgresDBConfig.Name, postgresDBConfig.Port, postgresDBConfig.SSLMode, postgresDBConfig.Timezone,
	)

	// postgres with sqlx
	postgresDBSqlx := InitializePostgresqlDatabaseSqlx(ctx, dsn)

	return DBCollection{
		PostgresDBSqlx: postgresDBSqlx,
	}
}
