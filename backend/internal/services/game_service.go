package services

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/111zxc/cocktaildle/backend/internal/db"
	"github.com/111zxc/cocktaildle/backend/internal/models"
)

type CocktailDetailAPIResponse struct {
	Drinks []Cocktail `json:"drinks"`
}

type Cocktail struct {
	IDDrink     string   `json:"idDrink"`
	StrDrink    string   `json:"strDrink"`
	StrAlcohol  string   `json:"strAlcoholic"`
	StrGlass    string   `json:"strGlass"`
	StrCategory string   `json:"strCategory"`
	Ingredients []string `json:"ingredients"`
}

func GetRandomCocktail() (*Cocktail, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := httpClient.Get("https://www.thecocktaildb.com/api/json/v1/1/random.php")
	if err != nil {
		fmt.Println("error get")
		fmt.Println(err)
		fmt.Println(resp)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return nil, errors.New("failed to fetch cocktail")
	}

	var cocktailResp CocktailDetailAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&cocktailResp); err != nil {
		return nil, err
	}

	if len(cocktailResp.Drinks) == 0 {
		return nil, errors.New("no cocktails found")
	}

	return &cocktailResp.Drinks[0], nil
}

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

func (s *GameService) SubmitGuess(userID string, guess string) (*models.Guess, map[string]bool, error) {
	dailyGame, err := s.GetDailyGame()
	if err != nil {
		return nil, nil, err
	}

	var gameAttempt models.GameAttempt
	if err := db.DB.Where("user_id = ? AND game_id = ?", userID, dailyGame.ID).First(&gameAttempt).Error; err != nil {
		return nil, nil, errors.New("игра для пользователя не найдена")
	}

	if gameAttempt.Correct {
		return nil, nil, errors.New("пользователь уже угадал коктейль")
	}

	correctCocktail, err := s.GetCocktailByID(dailyGame.CocktailID)
	if err != nil {
		return nil, nil, err
	}

	guessedCocktail, err := s.GetCocktailByID(guess)
	if err != nil {
		return nil, nil, err
	}

	isCorrect := guess == dailyGame.CocktailID

	categoryMatch := correctCocktail.StrCategory == guessedCocktail.StrCategory
	alcoholMatch := correctCocktail.StrAlcohol == guessedCocktail.StrAlcohol

	newGuess := models.Guess{
		GameAttemptID: gameAttempt.ID,
		Guess:         guess,
		IsCorrect:     isCorrect,
	}

	if err := db.DB.Create(&newGuess).Error; err != nil {
		return nil, nil, err
	}

	gameAttempt.GuessesMade++
	gameAttempt.LastGuessTime = time.Now()
	if isCorrect {
		gameAttempt.Correct = true
		dailyGame.PlayersWon++
	}

	if gameAttempt.GuessesMade > 0 {
		dailyGame.AvgGuesses = ((dailyGame.AvgGuesses * float64(dailyGame.PlayersPlayed)) + float64(gameAttempt.GuessesMade)) / float64(dailyGame.PlayersPlayed)
	}

	if gameAttempt.GuessesMade == 0 {
		dailyGame.PlayersPlayed++
	}

	if err := db.DB.Save(&gameAttempt).Error; err != nil {
		return nil, nil, err
	}
	if err := db.DB.Save(&dailyGame).Error; err != nil {
		return nil, nil, err
	}

	return &newGuess, map[string]bool{
		"correct":       isCorrect,
		"categoryMatch": categoryMatch,
		"alcoholMatch":  alcoholMatch,
	}, nil
}

func (s *GameService) GetCocktailByID(cocktailID string) (*Cocktail, error) {
	url := fmt.Sprintf("https://www.thecocktaildb.com/api/json/v1/1/lookup.php?i=%s", cocktailID)
	httpClient := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch cocktail by ID")
	}

	var cocktailResp CocktailDetailAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&cocktailResp); err != nil {
		return nil, err
	}

	if len(cocktailResp.Drinks) == 0 {
		return nil, errors.New("cocktail not found")
	}

	cocktail := cocktailResp.Drinks[0]
	cocktail.Ingredients = extractIngredients(cocktailResp.Drinks[0])

	return &cocktail, nil
}

func extractIngredients(Cocktail) []string {
	return make([]string, 5)
}
