package auth

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

func InitAuthHandlers(db *gorm.DB) {
	http.HandleFunc("/api/v1/login", loginHandler(db))
	http.HandleFunc("/api/v1/register", registerHandler(db))
}

type credentails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func validateParams(credentials *credentails) error {
	if credentials.Email == "" || credentials.Password == "" {
		return errors.New("Invalid request body")
	}
	return nil
}
