package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gitlab.com/back1ng1/messaggio-test/internal/entity"
	"gitlab.com/back1ng1/messaggio-test/internal/handlers"
	"gitlab.com/back1ng1/messaggio-test/internal/kafka"
	"gitlab.com/back1ng1/messaggio-test/internal/postgres"
	"gitlab.com/back1ng1/messaggio-test/internal/repository"
	"gitlab.com/back1ng1/messaggio-test/internal/usecase"
)

func main() {
	ctx := context.Background()
	done := make(chan struct{})

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
		Host:     os.Getenv("DB_HOST"),
		Port:     uint16(port),
		Database: os.Getenv("DB_DATABASE"),
	}

	pool, err := postgres.Setup(ctx, pgopt)
	if err != nil {
		log.Fatalf("failed create postgres instance: %s", err)
	}

	repo := repository.New(pool)
	broker := kafka.NewConn()

	uc := usecase.New(repo, broker)

	e := handlers.New(uc)
	go func() {
		e.Logger.Fatal(e.Start(":1323"))
	}()
	go func() {
		for msg := range kafka.NewReader() {
			var message entity.Message
			err := json.Unmarshal(msg.Value, &message)
			if err != nil {
				log.Fatal(err)
			}

			message, err = uc.ProcessMessage(message)
			if err != nil {
				fmt.Println(err)
			}
		}

	}()

	<-done
}
