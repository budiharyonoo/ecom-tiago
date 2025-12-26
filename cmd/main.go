package main

import (
	"database/sql"
	env "github/budiharyonoo/ecom-tiago/internal/utils"
	"log/slog"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbDsn := env.GetString("DB_DSN", "")

	db, err := sql.Open(env.GetString("DB_DRIVER", "mysql"), dbDsn)
	if err != nil {
		slog.Error("cannot connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	db.SetConnMaxLifetime(time.Minute * 1)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		slog.Error("cannot ping to the database", "error", err)
		os.Exit(1)
	}

	slog.Info("successfuly connected to the database")

	cfg := config{
		addr: ":3000",
		db: dbConfig{
			dsn: dbDsn,
		},
	}

	api := app{
		config: cfg,
		db:     db,
	}

	// run the server
	if err := api.run(api.mount()); err != nil {
		slog.Error("failed to start the server", "error", err)
		os.Exit(1)
	}
}
