package main

import (
	"banner-service/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

func main() {
	ctx := context.Background()
	cfg := config.NewConfig()
	conn, err := pgx.Connect(ctx, cfg.PgURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	if err := conn.Ping(ctx); err != nil {
		fmt.Print(err)
	}
	defer func() {
		if err := conn.Close(ctx); err != nil {
			log.Fatal(err)
		}
	}()

}
