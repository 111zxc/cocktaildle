package services

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/111zxc/cocktaildle/backend/internal/db"
	"github.com/111zxc/cocktaildle/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(username, email, password string) (*models.User, error) {
	var existingUser models.User
	if err := db.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	if err := db.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, errors.New("пользователь с таким username уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ошибка при хешировании пароля: %v", err)
		return nil, err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := db.DB.Create(user).Error; err != nil {
		log.Printf("ошибка при создании пользователя: %v", err)
		return nil, err
	}

	return user, nil
}

func (s *UserService) AuthenticateUser(email, password string) (string, error) {
	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("пользователь не найден")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("неверный пароль")
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		log.Printf("ошибка при генерации JWT токена: %v", err)
		return "", err
	}

	return token, nil
}

func (s *UserService) UpdateUserByID(userID string, newUser *models.User) (*models.User, error) {
	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, errors.New("пользователь не найден")
	}

	user.Username = newUser.Username
	user.Email = newUser.Email

	if newUser.PasswordHash != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("ошибка при хешировании пароля: %v", err)
			return nil, err
		}
		user.PasswordHash = string(hashedPassword)
	}

	if err := db.DB.Save(&user).Error; err != nil {
		log.Printf("ошибка при обновлении пользователя: %v", err)
		return nil, err
	}

	return &user, nil
}

func generateJWT(userID string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("отсутствует секретный ключ для JWT")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ValidateJWT(tokenString string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("отсутствует секретный ключ для JWT")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи токена")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("неверный формат user_id в токене")
		}
		return userID, nil
	}

	return "", errors.New("недействительный токен")
}
