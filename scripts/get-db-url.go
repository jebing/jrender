package main

import (
	"fmt"
	"log"
	"net/url"

	"revonoir.com/jrender/conns/configs"
)

func main() {
	configManager, err := configs.NewConfigManager("config")
	if err != nil {
		log.Fatalf("Failed to initialize config manager: %v", err)
	}

	config, err := configManager.GetConfig()
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	// Generate postgres URL for golang-migrate CLI with proper URL encoding
	dbConfig := config.Database
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(dbConfig.User),
		url.QueryEscape(dbConfig.Password),
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Dbname,
		dbConfig.Sslmode,
	)

	fmt.Print(dbURL)
}
