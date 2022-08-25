package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"joczkowski.com/room_keeper/auth"
	"joczkowski.com/room_keeper/users"
)

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println(os.Getenv("DB_URL"))
	dbPool, err := pgxpool.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer dbPool.Close()

	firstArg := os.Args[1]

	if firstArg == "automigrate" {
		fmt.Println("Automigrate")
	} else if firstArg == "server" {
		fmt.Println("Running server on port 8080...")

		users.InitUsersHandlers(dbPool)

		auth.InitAuthHandlers(dbPool)

		logPath := "development.log"

		openLogFile(logPath)

		httpPort := 8080

		http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
		// err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), logRequest(http.DefaultServeMux))
		// if err != nil {
		// 	log.Fatal(err)
		// }
	} else {
		fmt.Println("Invalid command")
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func openLogFile(logfile string) {
	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

		if err != nil {
			log.Fatal("OpenLogfile: os.OpenFile:", err)
		}

		log.SetOutput(lf)

		defer lf.Close()
	}
}
