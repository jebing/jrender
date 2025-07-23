package databases

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"revonoir.com/jbilling/conns/configs"
)

const SOURCE_DIR string = "file://resources/migrations"

func Migrate(ctx context.Context, config configs.Configuration) error {
	dbResourceString := GenerateDbString(config)

	db, err := sql.Open("postgres", dbResourceString)
	if err != nil {
		slog.Error("failed to get database connection for the migration", "error", err)
		return err
	}

	defer db.Close()

	dbConfig := config.Database

	if driver, err := postgres.WithInstance(db, &postgres.Config{}); err != nil {
		slog.Error("failed to build the database driver for migration", "error", err)
		return err
	} else if m, err := migrate.NewWithDatabaseInstance(
		SOURCE_DIR,
		dbConfig.Dbname,
		driver,
	); err != nil {
		slog.Error("failed creating database migration object", "error", err)
		return err
	} else if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			slog.Info("no migration to be done")
		} else {
			slog.Error("failed to migrate the database", "error", err)
			return err
		}
	}

	return nil
}
