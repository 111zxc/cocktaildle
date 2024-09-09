package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/111zxc/cocktaildle/backend/internal/services"
	"github.com/gorilla/mux"
)

type GameHandler struct {
	gameService *services.GameService
}

func NewGameHandler(gameService *services.GameService) *GameHandler {
	return &GameHandler{gameService: gameService}
}

func (h *GameHandler) StartGameHandler(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userID"]

	gameAttempt, err := h.gameService.StartGameForUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gameAttempt)
}

func (h *GameHandler) SubmitGuessHandler(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userID"]

	var input struct {
		Guess string `json:"guess"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	guess, err := h.gameService.SubmitGuess(userID, input.Guess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(guess)
}

func (h *GameHandler) GetDailyGameHandler(w http.ResponseWriter, r *http.Request) {
	dailyGame, err := h.gameService.GetDailyGame()
	if err != nil {
		dailyGame, _ = h.gameService.CreateDailyGame("asd")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dailyGame)
}
