package user_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jawahars16/gocloak/internal/user"
	"go.uber.org/mock/gomock"
)

func Test_AddUser_Success(t *testing.T) {
	reqBody := `{
		"first_name": "First Name",
		"last_name": "Last Name",
		"email": "Email",
		"password": "Password"
	}`
	ctrl := gomock.NewController(t)
	userManager := user.NewMockuserManager(ctrl)
	userManager.EXPECT().Add(user.User{
		FirstName: "First Name",
		LastName:  "Last Name",
		Email:     "Email",
		Password:  "Password",
	}).Times(1)

	handler := user.NewHandler(userManager)
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user", strings.NewReader(reqBody))
	handler.AddUser(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("%d not expected", rr.Code)
	}
}

func Test_AddUser_Badrequest(t *testing.T) {
	reqBody := `{
	}`
	ctrl := gomock.NewController(t)
	userManager := user.NewMockuserManager(ctrl)
	userManager.EXPECT().Add(user.User{}).Times(0)

	handler := user.NewHandler(userManager)
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user", strings.NewReader(reqBody))
	handler.AddUser(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("%d not expected", rr.Code)
	}
}
