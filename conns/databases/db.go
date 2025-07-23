package databases

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"revonoir.com/jbilling/conns/configs"
)

type Database struct {
	Client *pgxpool.Pool
}

func NewDatabase(ctx context.Context, config configs.Configuration) (*Database, error) {
	// generate pgxpool config
	pgxConfig, err := FormatConfig(config)
	if err != nil {
		slog.Error("Unable to generate pgxpool config", "error", err)
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close()
		slog.Error("Failed to ping the database", "error", err)
		return nil, err
	}

	db := &Database{
		Client: conn,
	}

	return db, nil
}

// Close closes the given *sql.DB connection.
// It logs information about the closing operation and reports any errors encountered.
func Close(db *pgxpool.Pool) {
	slog.Info("Closing database connection")
	db.Close()
}
