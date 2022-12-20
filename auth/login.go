package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var users = map[string]string{
	"user@example.com": "password",
}

var jwtKey = []byte("my_secret_key")

type claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func loginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentails credentails

		switch r.Method {
		case http.MethodPost:
			err := json.NewDecoder(r.Body).Decode(&credentails)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			expectedPassword, ok := users[credentails.Email]

			if !ok || expectedPassword != credentails.Password {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			expirationTime := time.Now().Add(5 * time.Minute)
			claims := &claims{
				Email: credentails.Email,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			})

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
