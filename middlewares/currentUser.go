package middlewares

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type CurrentUser struct {
	Email string
}

type AutheticatedHandler func(http.ResponseWriter, *http.Request, *pgxpool.Pool, *CurrentUser)

type EnsureAuth struct {
	handler AutheticatedHandler
	dbPool  *pgxpool.Pool
}

func (ea EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser{}

	ea.dbPool.QueryRow(context.Background(), "SELECT 'user@example.com' as email", r.Context().Value("user_id")).Scan(&user.Email)

	ea.handler(w, r, ea.dbPool, &user)
}

func NewEnsureAuth(handlerToWrap AutheticatedHandler, dbPool *pgxpool.Pool) *EnsureAuth {
	return &EnsureAuth{handlerToWrap, dbPool}
}