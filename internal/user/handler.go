package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type userManager interface {
	Add(user User) error
}

type handler struct {
	userManager userManager
}

func NewHandler(userManager userManager) handler {
	return handler{
		userManager: userManager,
	}
}

func (h *handler) AddUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		slog.Error("JSON decoder err", err)
	}

	if !validateUser(user) {
		slog.Error("User model invalid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.userManager.Add(user)
	if err != nil {
		slog.Error("Error adding user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func validateUser(user User) bool {
	return user.FirstName != "" &&
		user.LastName != "" &&
		user.Email != "" &&
		user.Password != ""
}
