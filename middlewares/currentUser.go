package middlewares

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
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

type AutheticatedHandler func(http.ResponseWriter, *http.Request, *pgxpool.Pool, *CurrentUser)

type EnsureAuth struct {
	handler AutheticatedHandler
	dbPool  *pgxpool.Pool
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

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
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

	ea.dbPool.QueryRow(context.Background(), "SELECT users.id, users.email, users.name  FROM (SELECT 1 as id,  'user@example.com' as email, 'Jakub Oczkowski' as name) as users WHERE users.email = $1", 1).Scan(&currentUser.Id, &currentUser.Email, &currentUser.Name)

	ea.handler(w, r, ea.dbPool, &currentUser)
}

func NewEnsureAuth(handlerToWrap AutheticatedHandler, dbPool *pgxpool.Pool) *EnsureAuth {
	return &EnsureAuth{handlerToWrap, dbPool}
}
