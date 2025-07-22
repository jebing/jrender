package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"revonoir.com/jform/conns/configs"
	"revonoir.com/jform/conns/databases"
	"revonoir.com/jform/webapp"
)

const (
	logFilePath = "/var/log/APP/jform"
	logFileName = "jform.log"
)

func main() {
	initLogger()

	slog.Info("GOMAXPROCS is set", "max_process", runtime.GOMAXPROCS(0))
	configManager, err := configs.NewConfigManager("config")
	if err != nil {
		slog.Error("failed to create config manager", "error", err)
		panic(err)
	}

	config, err := configManager.GetConfig()
	if err != nil {
		slog.Error("failed to get config", "error", err)
		panic(err)
	}

	db, err := databases.NewDatabase(context.Background(), config)
	if err != nil {
		slog.Error("failed to create database", "error", err)
		panic(err)
	}

	defer databases.Close(db.Client)

	r := initWebServer()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	webapp.Run(r)
}

func initWebServer() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*", "Accept", "Accept-Encoding", "User-Agent", "Host", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	return router
}

func initLogger() {
	var writers []io.Writer
	writers = append(writers, os.Stdout)

	// Configuring log file
	logFile, err := os.OpenFile(logFilePath+"/"+logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		slog.Error("Failed to create the log file", "error", err)
	} else {
		defer logFile.Close()
		writers = append(writers, logFile)
	}

	logHandler := slog.NewTextHandler(io.MultiWriter(writers...), &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	fileLogger := slog.New(logHandler)

	slog.SetDefault(fileLogger)
}
