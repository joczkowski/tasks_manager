package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type CurrentUser struct {
	Email string
}

type Handler func(dbPool *pgxpool.Pool, currentUser CurrentUser) http.HandlerFunc

func CurrentUserMiddleware(dbPool *pgxpool.Pool, handler Handler) http.HandlerFunc {
	var currentUser CurrentUser

	err := dbPool.QueryRow(context.Background(), "select 'user@example.com' as email").Scan(&currentUser.Email)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	return handler(dbPool, currentUser)
}
