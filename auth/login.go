package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
)

var users = map[string]string{
	"user@example.com": "password",
}

var jwtKey = []byte("my_secret_key")

type credentails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func InitAuthHandlers(dbPool *pgxpool.Pool) {
	http.HandleFunc("/v1/login", loginHandler(dbPool))
}

func loginHandler(dbPoll *pgxpool.Pool) http.HandlerFunc {
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
			claims := &Claims{
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
