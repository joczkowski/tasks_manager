package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"joczkowski.com/room_keeper/err_helpers"
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

func checkIfUserExist(credentials *credentails, db *gorm.DB) error {
	result := db.First(&User{}, "email = ?", credentials.Email)

	if result.RowsAffected > 0 {
		return errors.New("User already exist")
	}

	return nil
}

func registerHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var credentails credentails

			err := json.NewDecoder(r.Body).Decode(&credentails)
			err_helpers.HandleWebErr(w, err, http.StatusBadRequest)

			err = validateParams(&credentails)
			err_helpers.HandleWebErr(w, err, http.StatusBadRequest)

			err = checkIfUserExist(&credentails, db)
			err_helpers.HandleWebErr(w, err, http.StatusBadRequest)

			hashedPassword, err := HashPassword(credentails.Password)
			err_helpers.HandleWebErr(w, err, http.StatusBadRequest)

			result := db.Create(&User{
				Email:          credentails.Email,
				HashedPassword: hashedPassword,
			})
			err_helpers.HandleWebErr(w, result.Error, http.StatusBadRequest)

			w.WriteHeader(http.StatusCreated)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
