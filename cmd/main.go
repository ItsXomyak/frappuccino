package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"frappuccino/config"
	"frappuccino/internal/server"
	"frappuccino/internal/svc"

	repo "frappuccino/internal/repo"

	_ "github.com/lib/pq"
)

var (
	port  int
	dbURL string
)

func main() {
	config.LoadEnv()

	flag.IntVar(&port, "port", config.GetEnvInt("PORT", 9090), "Port number")
	flag.StringVar(&dbURL, "db", os.Getenv("DATABASE_URL"), "Database connection URL")
	flag.Parse()

	if port < 0 || port > 65535 {
		log.Fatal("Invalid port number")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection test failed: %v", err)
	}

	fmt.Println("Connected to database successfully")

	container := repo.New(db)

	service := svc.NewSvc(container)

	handler := server.New(service)

	srv := server.NewServer(strconv.Itoa(port), *handler)
	srv.Start()
}
