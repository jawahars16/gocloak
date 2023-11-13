package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jawahars16/gocloak/internal/user"
)

type authResponse struct {
	AccessToken string `json:"access_token"`
}

type userManager interface {
	GenerateToken(email string, password string) (string, error)
}

type handler struct {
	userManager userManager
}

func NewHandler(userManager userManager) handler {
	return handler{
		userManager: userManager,
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var user user.User
	json.NewDecoder(r.Body).Decode(&user)

	token, err := h.userManager.GenerateToken(user.Email, user.Password)
	if errors.Is(err, ErrUserNotFound) {
		slog.Error("User not found", slog.String("email", user.Email))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if errors.Is(err, ErrMismatchedHashAndPassword) {
		slog.Error("Mismatched hash and password", slog.String("email", user.Email))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err != nil {
		slog.Error("Error generating token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authResponse{AccessToken: token})
}
