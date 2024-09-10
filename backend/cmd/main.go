package main

import (
	"log"
	"net/http"

	"github.com/111zxc/cocktaildle/backend/internal/db"
	"github.com/111zxc/cocktaildle/backend/internal/handlers"
	"github.com/111zxc/cocktaildle/backend/internal/middleware"
	"github.com/111zxc/cocktaildle/backend/internal/services"
	"github.com/gorilla/mux"
)

func main() {
	db.ConnectDatabase()

	gameService := services.NewGameService()
	gameHandler := handlers.NewGameHandler(gameService)

	userService := services.NewUserService()
	userHandler := handlers.NewUserHandler(userService)

	r := mux.NewRouter()

	r.Handle("/api/game/start", middleware.JWTAuthMiddleware(http.HandlerFunc(gameHandler.StartGameHandler))).Methods("POST")
	r.Handle("/api/game/guess", middleware.JWTAuthMiddleware(http.HandlerFunc(gameHandler.SubmitGuessHandler))).Methods("POST")
	r.Handle("/api/user/{id}", middleware.JWTAuthMiddleware(http.HandlerFunc(userHandler.UpdateUserHandler))).Methods("PUT")

	r.HandleFunc("/api/register", userHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", userHandler.LoginHandler).Methods("POST")
	r.HandleFunc("/api/game/daily", gameHandler.GetDailyGameHandler).Methods("GET")

	http.Handle("/", r)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
