package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/wagnerbrenner/study-go-and-react/internal/api"
	"github.com/wagnerbrenner/study-go-and-react/internal/store/pgstore"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("WSRS_DATABASE_USER"),
		os.Getenv("WSRS_DATABASE_PASSWORD"),
		os.Getenv("WSRS_DATABASE_HOST"),
		os.Getenv("WSRS_DATABASE_PORT"),
		os.Getenv("WSRS_DATABASE_NAME"),
	))
	if err != nil {
		panic(fmt.Sprintf("Error creating connection pool: %v", err))
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(fmt.Sprintf("Error pinging database: %v", err))
	}

	handler := api.NewHandler(pgstore.New(pool))

	go func() {
		if err := http.ListenAndServe(":8080", handler); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(fmt.Sprintf("Error starting server: %v", err))
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	fmt.Println("Server started. Press Ctrl+C to stop.")
	<-quit
}
