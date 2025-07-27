package webapp

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type App struct {
	server    *http.Server
	waitGroup *sync.WaitGroup
	ch        *chan os.Signal
}

func Run(r *chi.Mux) {
	app := New(r)
	app.ListenAndServe()
}

func New(r *chi.Mux) *App {
	s := &http.Server{
		Addr:         ":9200",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	wg := &sync.WaitGroup{}
	ch := make(chan os.Signal, 1)

	wg.Add(1)
	go func() {
		// wait for interrupt signal
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		sig := <-ch
		slog.Info("Received signal, initiating server shutdown", "signal", sig)

		// gracefully shutdown the server with a timeout of 10 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			slog.Info("Error while waiting for server to shutdown", "error", err)
		}

		slog.Info("Server shutdown complete")
		wg.Done()
	}()

	return &App{
		server:    s,
		waitGroup: wg,
		ch:        &ch,
	}
}

func (a *App) RegisterOnShutdown(f func()) *App {
	a.waitGroup.Add(1)
	a.server.RegisterOnShutdown(func() {
		f()
		a.waitGroup.Done()
	})
	return a
}

func (a *App) ListenAndServe() {
	slog.Info("Starting application and listening to", "address", a.server.Addr)
	if err := a.server.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("Error while starting the server", "error", err)
	}

	a.waitGroup.Wait()
	slog.Info("Done waiting for server to shutdown")
}
