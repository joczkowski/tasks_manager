package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"joczkowski.com/room_keeper/err_helpers"
)

var users = map[string]string{
	"user@example.com": "password",
}

var jwtKey = []byte("my_secret_key")

type claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func findUserByEamilAndPassword(credentials *credentails, db *gorm.DB) (User, error) {
	var user User
	result := db.First(&user, "email = ?", credentials.Email)

	if result.RowsAffected == 0 {
		return User{}, errors.New("User not found")
	}

	if !checkPasswordHash(credentials.Password, user.HashedPassword) {
		return User{}, errors.New("Invalid password")
	}

	return user, nil
}

func createToken(user *User) (time.Time, string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return expirationTime, "", err
	}

	return expirationTime, tokenString, nil
}

func loginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var credentails credentails

			_ = json.NewDecoder(r.Body).Decode(&credentails)

			err := validateParams(&credentails)
			err_helpers.HandleWebErr(w, err, http.StatusBadRequest)

			user, err := findUserByEamilAndPassword(&credentails, db)
			err_helpers.HandleWebErr(w, err, http.StatusBadRequest)

			expireTime, tokenString, err := createToken(&user)
			err_helpers.HandleWebErr(w, err, http.StatusBadRequest)

			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expireTime,
			})
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
