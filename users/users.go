package users

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"joczkowski.com/room_keeper/middlewares"
)

type me struct {
	Email string `json:"email"`
}

type data struct {
	Me me `json:"user"`
}

type jsonResponse struct {
	Status string `json:"status"`
	Data   data   `json:"data"`
}

func InitUsersHandlers(dbPool *pgxpool.Pool) {
	http.Handle("/v1/users/me", middlewares.CurrentUserMiddleware(dbPool, meHandler))
}

func meHandler(dbPool *pgxpool.Pool, currentUser middlewares.CurrentUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonResponse := jsonResponse{
			Status: "ok",
			Data:   data{Me: me{Email: currentUser.Email}},
		}

		jsonData, err := json.Marshal(jsonResponse)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}
