package auth_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jawahars16/gocloak/internal/auth"
	"go.uber.org/mock/gomock"
)

func Test_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockManager := auth.NewMockuserManager(ctrl)

	t.Run("should return error if user not found", func(t *testing.T) {
		userBody := `
		{
			"email": "test@test.com",
			"password": "password"
		}
		`
		mockManager.EXPECT().GenerateToken("test@test.com", "password").Return("", auth.ErrUserNotFound)
		hanlder := auth.NewHandler(mockManager)

		rr := httptest.NewRecorder()
		rew := httptest.NewRequest("POST", "/login", strings.NewReader(userBody))
		hanlder.Login(rr, rew)

		if rr.Code != 401 {
			t.Errorf("Expected 401, got %v", rr.Code)
		}
	})

	t.Run("should return error if password is incorrect", func(t *testing.T) {
		userBody := `
		{
			"email": "test@test.com",
			"password": "password"
		}
		`
		mockManager.EXPECT().GenerateToken("test@test.com", "password").Return("", auth.ErrMismatchedHashAndPassword)
		hanlder := auth.NewHandler(mockManager)

		rr := httptest.NewRecorder()
		rew := httptest.NewRequest("POST", "/login", strings.NewReader(userBody))
		hanlder.Login(rr, rew)

		if rr.Code != 401 {
			t.Errorf("Expected 401, got %v", rr.Code)
		}
	})

	t.Run("should return nil if password is correct", func(t *testing.T) {
		userBody := `
		{
			"email": "test@test.com",
			"password": "password"
		}
		`
		mockManager.EXPECT().GenerateToken("test@test.com", "password").Return("token", nil)
		hanlder := auth.NewHandler(mockManager)

		rr := httptest.NewRecorder()
		rew := httptest.NewRequest("POST", "/login", strings.NewReader(userBody))
		hanlder.Login(rr, rew)

		if rr.Code != 200 {
			t.Errorf("Expected 200, got %v", rr.Code)
		}

		if rr.Body.String() != "token" {
			t.Errorf("Expected token, got %v", rr.Body.String())
		}
	})
}
