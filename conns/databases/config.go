package databases

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"revonoir.com/jform/conns/configs"
)

type DBEnv struct {
	DBUser           string `envName:"DB_USER"`
	DBPassword       string `envName:"DB_PASSWORD"`
	DBHost           string `envName:"DB_HOST"`
	DBPort           int    `envName:"DB_PORT" defaultValue:"5432"`
	DBName           string `envName:"DB_NAME"`
	DBMaxConnections int32  `envName:"NV_DB_MAX_CONNS" defaultValue:"10"`
}

// convert the config into the postgres URL string config
// example: user=test_user password=<PASSWORD> host=test_host port=1234 dbname=test_db sslmode=disable
func GenerateDbString(config configs.Configuration) string {
	dbConfig := config.Database
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s pool_max_conns=%d sslmode=disable",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Dbname,
		dbConfig.MaxConns,
	)
}

func FormatConfig(config configs.Configuration) (*pgxpool.Config, error) {
	pgxConfigString := GenerateDbString(config)

	if pgxConfig, err := pgxpool.ParseConfig(pgxConfigString); err != nil {
		return nil, err
	} else {
		// Setting up the pool
		pgxConfig.MaxConns = int32(config.Database.MaxConns)
		pgxConfig.MaxConnLifetime = 5 * time.Minute

		return pgxConfig, nil
	}
}
