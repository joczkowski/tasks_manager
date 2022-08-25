package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"joczkowski.com/room_keeper/middlewares"
)

type me struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Id    int    `json:"id"`
}

type data struct {
	Me me `json:"user"`
}

type jsonResponse struct {
	Status string `json:"status"`
	Data   data   `json:"data"`
}

func InitUsersHandlers(dbPool *pgxpool.Pool) {
	http.Handle("/v1/me", middlewares.NewEnsureAuth(meHandler, dbPool))
}

func meHandler(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool, currentUser *middlewares.CurrentUser) {
	fmt.Println(currentUser)
	jsonResponse := jsonResponse{
		Status: "ok",
		Data:   data{Me: me{Email: currentUser.Email, Name: currentUser.Name, Id: currentUser.Id}},
	}

	jsonData, err := json.Marshal(jsonResponse)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
