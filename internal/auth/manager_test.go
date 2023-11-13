package auth_test

import (
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jawahars16/gocloak/config"
	"github.com/jawahars16/gocloak/internal/auth"
	"github.com/jawahars16/gocloak/internal/user"
	gomock "go.uber.org/mock/gomock"
)

func Test_GenerateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDB := auth.NewMockdb(ctrl)
	mockvalidator := auth.NewMockvalidator(ctrl)
	mocktokenGenerator := auth.NewMocktokenGenerator(ctrl)
	config := config.AuthConfig{JWTSecret: "secret"}

	t.Run("should return error if user not found", func(t *testing.T) {
		mockDB.EXPECT().First(gomock.Any(), gomock.Any()).Return(errors.New("user not found"))
		manager := auth.NewManager(mockDB, mockvalidator, mocktokenGenerator, config)
		_, err := manager.GenerateToken("test@test.com", "password")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("should return error if password is incorrect", func(t *testing.T) {
		mockvalidator.EXPECT().CompareHashAndPassword(gomock.Any(), gomock.Any()).Return(errors.New("password incorrect"))
		mockDB.EXPECT().First(gomock.Any(), gomock.Any()).Return(nil)
		manager := auth.NewManager(mockDB, mockvalidator, mocktokenGenerator, config)
		_, err := manager.GenerateToken("test@test.com", "password")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("should return nil if password is correct", func(t *testing.T) {
		userData := &user.User{}

		mockvalidator.EXPECT().CompareHashAndPassword(gomock.Any(), gomock.Any()).Return(nil)
		mockDB.EXPECT().First(userData, "email = ?", "test@test.com").Return(nil)
		mocktokenGenerator.EXPECT().NewWithClaims(jwt.SigningMethodHS256, gomock.Any()).Return(&jwt.Token{
			Method: jwt.SigningMethodHS256,
		})

		manager := auth.NewManager(mockDB, mockvalidator, mocktokenGenerator, config)
		_, err := manager.GenerateToken("test@test.com", "password")
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})
}
