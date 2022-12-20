package auth

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id             int `gorm:"primaryKey"`
	Email          string
	HashedPassword string
}

func registerHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var credentails credentails

			err := json.NewDecoder(r.Body).Decode(&credentails)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			result := db.Create(&User{
				Email:          credentails.Email,
				HashedPassword: credentails.Password,
			})

			if result.Error != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			return
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
