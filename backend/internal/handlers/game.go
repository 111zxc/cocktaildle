package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/111zxc/cocktaildle/backend/internal/middleware"
	"github.com/111zxc/cocktaildle/backend/internal/models"
	"github.com/111zxc/cocktaildle/backend/internal/services"
)

type GameHandler struct {
	gameService *services.GameService
}

func NewGameHandler(gameService *services.GameService) *GameHandler {
	return &GameHandler{gameService: gameService}
}

func (h *GameHandler) StartGameHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем user_id из контекста
	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok {
		http.Error(w, "Unable to retrieve user_id", http.StatusUnauthorized)
		return
	}

	// Вызываем логику для начала игры
	gameAttempt, err := h.gameService.StartGameForUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gameAttempt)
}

func (h *GameHandler) SubmitGuessHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем user_id из контекста
	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok {
		http.Error(w, "Unable to retrieve user_id", http.StatusUnauthorized)
		return
	}

	// Читаем догадку пользователя из тела запроса
	var input struct {
		Guess string `json:"guess"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Вызываем логику для обработки догадки
	guess, details, err := h.gameService.SubmitGuess(userID, input.Guess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Guess   *models.Guess   `json:"guess"`
		Details map[string]bool `json:"details"`
	}{
		Guess:   guess,
		Details: details,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *GameHandler) GetDailyGameHandler(w http.ResponseWriter, r *http.Request) {
	dailyGame, err := h.gameService.GetDailyGame()
	if err != nil {
		cocktail, err := services.GetRandomCocktail()
		if err != nil {
			http.Error(w, "Failed to fetch random cocktail", http.StatusInternalServerError)
			return
		}

		dailyGame, err = h.gameService.CreateDailyGame(cocktail.IDDrink)
		if err != nil {
			http.Error(w, "Failed to create daily game", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dailyGame)
}
