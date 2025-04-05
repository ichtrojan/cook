package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/ichtrojan/cook/config"
	"github.com/ichtrojan/cook/database"
	"github.com/ichtrojan/cook/queue"
	"github.com/ichtrojan/cook/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func monitoring() *http.Server {
	monitoring := asynqmon.New(asynqmon.Options{
		RootPath:     "/monitoring",
		RedisConnOpt: asynq.RedisClientOpt{Addr: ":" + config.RedisConfig.Port},
	})

	router := mux.NewRouter()
	router.PathPrefix(monitoring.RootPath()).Handler(monitoring)

	return &http.Server{
		Handler: monitoring,
		Addr:    ":6660",
	}
}

func main() {
	if err := config.LoadEnvironmentVariables(); err != nil {
		log.Fatal(err)
	}

	if err := database.ConnectMySQL(config.MySQLConfig); err != nil {
		log.Fatal(err)
	}

	if err := database.ConnectRedis(config.RedisConfig); err != nil {
		log.Fatal(err)
	}

	serverError := make(chan error, 1)

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	cookServer := &http.Server{
		Addr:    ":6666",
		Handler: routes.AllRoutes(),
	}

	if config.AppConfig.AsyncmonService == "true" {
		go func() {
			fmt.Println("Starting Monitoring server on :6660")
			if err := monitoring().ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				serverError <- fmt.Errorf("monitoring server error: %v", err)
			}
		}()
	}

	go func() {
		fmt.Println("Starting Cook server on :6666")

		if err := cookServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverError <- fmt.Errorf("cook server error: %v", err)
		}
	}()

	go func() {
		fmt.Println("Starting Asynq worker server")
		worker := asynq.NewServer(
			asynq.RedisClientOpt{Addr: database.Redis.Options().Addr},
			asynq.Config{
				Concurrency: 10,
			},
		)

		if err := worker.Run(queue.Register()); err != nil {
			serverError <- fmt.Errorf("asynq worker error: %v", err)
		}
	}()

	// Wait for an OS signal or a server error
	select {
	case <-stop:
		fmt.Println("Shutting down...")
	case err := <-serverError:
		log.Printf("server error: %v", err)
	}

	// Gracefully shutdown servers
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := monitoring().Shutdown(ctx); err != nil {
		log.Printf("Asyncmon server shutdown failed: %v", err)
	}

	if err := cookServer.Shutdown(ctx); err != nil {
		log.Printf("Cook server shutdown failed: %v", err)
	}

	if err := queue.Client.Close(); err != nil {
		log.Printf("Failed to close redis queue: %v", err)
	}

	fmt.Println("Servers have been stopped gracefully.")
}
