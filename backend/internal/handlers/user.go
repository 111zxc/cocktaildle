package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/111zxc/cocktaildle/backend/internal/models"
	"github.com/111zxc/cocktaildle/backend/internal/services"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		services.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	user, err := h.userService.CreateUser(input.Username, input.Email, input.Password)
	if err != nil {
		services.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    user,
	})
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		services.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	token, err := h.userService.AuthenticateUser(input.Email, input.Password)
	if err != nil {
		services.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token":   token,
	})
}

func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		services.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	updatedUser := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: input.Password,
	}

	user, err := h.userService.UpdateUserByID(userID, updatedUser)
	if err != nil {
		services.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    user,
	})
}
