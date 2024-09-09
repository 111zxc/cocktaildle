package services

import (
	"errors"
	"log"
	"time"

	"github.com/111zxc/cocktaildle/backend/internal/db"
	"github.com/111zxc/cocktaildle/backend/internal/models"
)

type GameService struct{}

func NewGameService() *GameService {
	return &GameService{}
}

func (s *GameService) GetDailyGame() (*models.DailyGame, error) {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	var dailyGame models.DailyGame
	if err := db.DB.Where("date = ?", today).First(&dailyGame).Error; err == nil {
		return &dailyGame, nil
	}
	return nil, errors.New("игры не существует")
}

func (s *GameService) CreateDailyGame(cocktailID string) (*models.DailyGame, error) {
	_, err := s.GetDailyGame()
	if err == nil {
		return nil, errors.New("игра уже существует")
	}

	dailyGame := models.DailyGame{
		CocktailID:    cocktailID,
		Date:          time.Now().UTC().Truncate(24 * time.Hour),
		PlayersPlayed: 0,
		PlayersWon:    0,
		AvgGuesses:    0.0,
	}

	if err := db.DB.Create(&dailyGame).Error; err != nil {
		log.Printf("ошибка при создании ежедневной игры: %v", err)
		return nil, err
	}

	return &dailyGame, nil
}

func (s *GameService) StartGameForUser(userID string) (*models.GameAttempt, error) {
	dailyGame, err := s.GetDailyGame()
	if err != nil {
		return nil, err
	}

	var existingAttempt models.GameAttempt
	if err := db.DB.Where("user_id = ? AND game_id = ?", userID, dailyGame.ID).First(&existingAttempt).Error; err == nil {

		return &existingAttempt, nil
	}

	newAttempt := models.GameAttempt{
		UserID:        userID,
		GameID:        dailyGame.ID,
		GuessesMade:   0,
		Correct:       false,
		LastGuessTime: time.Now(),
	}

	if err := db.DB.Create(&newAttempt).Error; err != nil {
		return nil, err
	}

	return &newAttempt, nil
}

func (s *GameService) SubmitGuess(userID string, guess string) (*models.Guess, error) {
	dailyGame, err := s.GetDailyGame()
	if err != nil {
		return nil, err
	}

	var gameAttempt models.GameAttempt
	if err := db.DB.Where("user_id = ? AND game_id = ?", userID, dailyGame.ID).First(&gameAttempt).Error; err != nil {
		return nil, errors.New("игра для пользователя не найдена")
	}

	if gameAttempt.Correct {
		return nil, errors.New("пользователь уже угадал коктейль")
	}

	isCorrect := guess == dailyGame.CocktailID

	newGuess := models.Guess{
		GameAttemptID: gameAttempt.ID,
		Guess:         guess,
		IsCorrect:     isCorrect,
	}

	if err := db.DB.Create(&newGuess).Error; err != nil {
		return nil, err
	}

	gameAttempt.GuessesMade++
	gameAttempt.LastGuessTime = time.Now()
	if isCorrect {
		gameAttempt.Correct = true
		dailyGame.PlayersWon++
	}

	if gameAttempt.GuessesMade > 0 {
		dailyGame.AvgGuesses = ((dailyGame.AvgGuesses * float64(dailyGame.PlayersPlayed)) + float64(gameAttempt.GuessesMade)) / float64(dailyGame.PlayersPlayed+1)
	}

	dailyGame.PlayersPlayed++

	if err := db.DB.Save(&gameAttempt).Error; err != nil {
		return nil, err
	}
	if err := db.DB.Save(&dailyGame).Error; err != nil {
		return nil, err
	}

	return &newGuess, nil
}
