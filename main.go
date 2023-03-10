package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"joczkowski.com/room_keeper/auth"
	"joczkowski.com/room_keeper/tasks"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupDataBase() *gorm.DB {
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	if dbUser == "" || dbName == "" || dbPassword == "" {
		panic("Missing database credentials")
	}

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s", dbUser, dbPassword, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func setupEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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

func main() {
	setupEnvVariables()
	db := setupDataBase()

	firstArg := os.Args[1]

	if firstArg == "migrate" {
		db.AutoMigrate(&User{})
		db.AutoMigrate(&Task{})
	} else if firstArg == "server" {
		fmt.Println("Running server on port 8000...")

		auth.InitAuthHandlers(db)
		tasks.InitTaskHandlers(db)

		if os.Getenv("PROJECT_ENV") == "development" {
			http.HandleFunc("/admin/users", func(w http.ResponseWriter, r *http.Request) {
				type userResult struct {
					ID    int
					Email string
				}

				var users []userResult
				db.Table("users").Find(&users)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				json.NewEncoder(w).Encode(users)
			})
		}

		http.ListenAndServe(fmt.Sprintf(":%d", 8000), logRequest(http.DefaultServeMux))
	} else {
		fmt.Println("Invalid command")
	}
}
