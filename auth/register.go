package auth

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id             int `gorm:"primaryKey"`
	Email          string
	HashedPassword string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
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

			if credentails.Email == "" || credentails.Password == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Empty params"})
				return
			}

			result := db.First(&User{}, "email = ?", credentails.Email)

			if result.RowsAffected > 0 {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid params"})
				return
			}

			hashedPassword, err := HashPassword(credentails.Password)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			result = db.Create(&User{
				Email:          credentails.Email,
				HashedPassword: hashedPassword,
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
