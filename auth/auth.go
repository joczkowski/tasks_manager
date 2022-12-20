package auth

import (
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
