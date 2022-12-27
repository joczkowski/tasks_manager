package tasks

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
	"joczkowski.com/room_keeper/err_helpers"
	"joczkowski.com/room_keeper/middlewares"
)

func InitTaskHandlers(db *gorm.DB) {
	http.Handle("/api/v1/tasks", middlewares.NewEnsureAuth(allTask, db))
	http.Handle("/api/v1/task", middlewares.NewEnsureAuth(createTask, db))
	http.Handle("/api/v1/task/", middlewares.NewEnsureAuth(idBaseActions, db))
}

func allTask(w http.ResponseWriter, r *http.Request, db *gorm.DB, currentUser *middlewares.CurrentUser) {
	switch r.Method {
	case "GET":
		type task struct {
			Id          int `gorm:"primaryKey"`
			Title       string
			Description string
			UserId      int
			CreatedAt   time.Time
			UpdatedAt   time.Time
		}

		var tasks []task

		db.Table("tasks").Where("user_id = ?", currentUser.Id).Find(&tasks)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(tasks)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func createTask(w http.ResponseWriter, r *http.Request, db *gorm.DB, currentUser *middlewares.CurrentUser) {
	switch r.Method {
	case "POST":
		type taskForm struct {
			Title       string
			Description string
			UserId      int
		}

		var task taskForm

		err := json.NewDecoder(r.Body).Decode(&task)
		err_helpers.HandleWebErr(w, err, http.StatusBadRequest)

		task.UserId = currentUser.Id

		result := db.Table("tasks").Create(&task)
		if result.Error != nil {
			err_helpers.HandleWebErr(w, result.Error, http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func idBaseActions(w http.ResponseWriter, r *http.Request, db *gorm.DB, currentUser *middlewares.CurrentUser) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/task/")

	switch r.Method {
	case "PATCH":
		type taskForm struct {
			Title       string
			Description string
		}

		var task taskForm

		err := json.NewDecoder(r.Body).Decode(&task)
		err_helpers.HandleWebErr(w, err, http.StatusBadRequest)

		result := db.Table("tasks").Where("id = ?", id).Updates(&task)

		err_helpers.HandleWebErr(w, result.Error, http.StatusBadRequest)
		if result.RowsAffected == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	case "DELETE":
		result := db.Exec("DELETE FROM tasks WHERE id = ?", id)
		err_helpers.HandleWebErr(w, result.Error, http.StatusBadRequest)
		if result.RowsAffected == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	case "GET":
		type taskModel struct {
			Id          int
			Title       string
			Description string
			CreatedAt   time.Time
			UpdatedAt   time.Time
		}

		var task taskModel

		result := db.Table("tasks").Where("id = ?", id).Find(&task)
		if result.RowsAffected == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(task)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
