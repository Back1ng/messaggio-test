package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gitlab.com/back1ng1/messaggio-test/internal/handlers"
	"gitlab.com/back1ng1/messaggio-test/internal/postgres"
	"gitlab.com/back1ng1/messaggio-test/internal/repository"
	"gitlab.com/back1ng1/messaggio-test/internal/usecase"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("failed parse DB_PORT: %s", err)
	}

	pgopt := postgres.SetupOptions{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		//Host:     os.Getenv("DB_HOST"),
		Host:     "127.0.0.1",
		Port:     uint16(port),
		Database: os.Getenv("DB_DATABASE"),
	}

	pool, err := postgres.Setup(ctx, pgopt)
	if err != nil {
		log.Fatalf("failed create postgres instance: %s", err)
	}

	repo := repository.New(pool)
	// todo add kafka to usecase
	uc := usecase.New(repo)

	e := handlers.New(uc)

	e.Logger.Fatal(e.Start(":1323"))

	// get http messages
	// store them into postgresql and send to kafka
	// after processing change status in postgresql with marking timestamp

	// http handler for showing this stats
}
