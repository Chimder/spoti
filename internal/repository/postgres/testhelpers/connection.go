package testhelpers

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func SetupContainerDB() *pgxpool.Pool {
	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:latest",
		tcpostgres.WithDatabase("test"),
		tcpostgres.WithUsername("user"),
		tcpostgres.WithPassword("password"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		panic(fmt.Sprintf("err start pgContainer: %v", err))
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("err pgContainer conn: %v", err))
	}

	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		panic(fmt.Sprintf("err to open sql: %v", err))
	}

	err = goose.Up(sqlDB, "../../../../sql/migrations/postgres")
	if err != nil {
		panic(fmt.Sprintf("err goose up: %v", err))
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		panic(fmt.Sprintf("err pgxpoolNew: %v", err))
	}

	return pool
}
