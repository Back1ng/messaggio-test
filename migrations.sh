#!/bin/bash
source .env
goose -dir migrations/ postgres "host=127.0.0.1 user=${DB_USERNAME} dbname=${DB_DATABASE} password=${DB_PASSWORD} port=${DB_PORT} sslmode=disable" up