package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/111zxc/cocktaildle/backend/internal/middleware"
	"github.com/111zxc/cocktaildle/backend/internal/services"
)

type GameHandler struct {
	gameService *services.GameService
}

func NewGameHandler(gameService *services.GameService) *GameHandler {
	return &GameHandler{gameService: gameService}
}

func (h *GameHandler) StartGameHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok {
		services.RespondWithError(w, http.StatusUnauthorized, "Unable to retrieve user_id")
		return
	}

	gameAttempt, err := h.gameService.StartGameForUser(userID)
	if err != nil {
		services.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"game":    gameAttempt,
	})
}

func (h *GameHandler) SubmitGuessHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok {
		services.RespondWithError(w, http.StatusUnauthorized, "Unable to retrieve user_id")
		return
	}

	var input struct {
		Guess string `json:"guess"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		services.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	guess, details, err := h.gameService.SubmitGuess(userID, input.Guess)
	if err != nil {
		services.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"guess":   guess,
		"details": details,
	})
}

func (h *GameHandler) GetDailyGameHandler(w http.ResponseWriter, r *http.Request) {
	dailyGame, err := h.gameService.GetDailyGame()
	if err != nil {
		cocktail, err := services.GetRandomCocktail()
		if err != nil {
			services.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch random cocktail")
			return
		}

		dailyGame, err = h.gameService.CreateDailyGame(cocktail.IDDrink)
		if err != nil {
			services.RespondWithError(w, http.StatusInternalServerError, "Failed to create daily game")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"game":    dailyGame,
	})
}
