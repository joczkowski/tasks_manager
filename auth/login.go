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

type postLoginAction struct {
	w http.ResponseWriter
	r *http.Request
	db *gorm.DB
	deseralizedParams credentails
	tokenString string
	expirationTime time.Time
}

func (self *postLoginAction) handle()  {
	if !self.deserializeParams() { return }
	if !self.validateParams() { return }
	if !self.auth() { return }

	http.SetCookie(self.w, &http.Cookie{
		Name:    "token",
		Value:   self.tokenString,
		Expires: self.expirationTime,
	})
}

func (self *postLoginAction) deserializeParams() bool {
	err := json.NewDecoder(self.r.Body).Decode(&self.deseralizedParams)
	if err != nil {
		self.w.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
}

func (self *postLoginAction) validateParams() bool {
	if self.deseralizedParams.Email == "" || self.deseralizedParams.Password == "" {
		self.w.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
}

func (self *postLoginAction) auth() bool {
	var user User

	result := self.db.First(user, "email = ?", self.deseralizedParams.Email)

	if result.RowsAffected == 0 {
		self.w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	
	if !CheckPasswordHash(self.deseralizedParams.Password, user.HashedPassword) {
		self.w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	self.expirationTime = time.Now().Add(5 * time.Minute)
	claims := &claims{
		Email: self.deseralizedParams.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: self.expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var err error
	self.tokenString, err = token.SignedString(jwtKey)
	if err != nil {
		self.w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}

func loginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			action := postLoginAction{w: w, r: r, db: db}
			action.handle()
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
