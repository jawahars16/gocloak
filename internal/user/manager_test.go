package user_test

import (
	"testing"

	"github.com/jawahars16/gocloak/internal/user"
	gomock "go.uber.org/mock/gomock"
)

func Test_AddUser(t *testing.T) {
	userData := user.User{FirstName: "firstname", LastName: "lastname", Email: "email", Password: "password"}
	ctrl := gomock.NewController(t)

	mockDB := user.NewMockdb(ctrl)
	mockCrypto := user.NewMockcrypto(ctrl)

	mockDB.EXPECT().Save(&user.User{FirstName: "firstname", LastName: "lastname", Email: "email", Password: "hashed-password"}).Return(nil).Times(1)
	mockCrypto.EXPECT().GenerateFromPassword(userData.Password).Return("hashed-password", nil).Times(1)

	manager := user.NewManager(mockDB, mockCrypto)
	err := manager.Add(userData)
	if err != nil {
		t.Error("Unexpected error", err)
	}
}
