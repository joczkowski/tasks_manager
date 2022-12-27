package middlewares

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"joczkowski.com/room_keeper/err_helpers"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var jwtKey = []byte("my_secret_key")

type CurrentUser struct {
	Id    int
	Email string
	Name  string
}

type AutheticatedHandler func(http.ResponseWriter, *http.Request, *gorm.DB, *CurrentUser)

type EnsureAuth struct {
	handler AutheticatedHandler
	db      *gorm.DB
}

func (ea EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")

	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tknStr := c.Value

	claims := &Claims{}

	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			err_helpers.HandleWebErr(w, err, http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var currentUser CurrentUser

	ea.db.Where("email = ?", claims.Email).First(&currentUser)

	ea.handler(w, r, ea.db, &currentUser)
}

func NewEnsureAuth(handlerToWrap AutheticatedHandler, db *gorm.DB) *EnsureAuth {
	return &EnsureAuth{handlerToWrap, db}
}
