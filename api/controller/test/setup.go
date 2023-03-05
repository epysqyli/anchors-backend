package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	routev1 "github.com/epysqyli/anchors-backend/api/route/v1"
	"github.com/epysqyli/anchors-backend/bootstrap"
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func assertEqual(exp any, got any, t *testing.T, info ...string) {
	if exp != got {
		if len(info) == 0 {
			t.Fatalf("\n\n-- expected: %v\n-- obtained: %v\n\n", exp, got)
		} else {
			t.Fatalf("\n\n%s\n-- expected: %v\n-- obtained: %v\n\n", info[0], exp, got)
		}
	}
}

func assertUnequal(value any, got any, t *testing.T, info ...string) {
	if value == got {
		if len(info) == 0 {
			t.Fatalf("\n\n-- value: %v\n-- obtained: %v\n\n", value, got)
		} else {
			t.Fatalf("\n\n%s\n-- value: %v\n-- obtained: %v\n\n", info[0], value, got)
		}
	}
}

func setup() (*gin.Engine, *gorm.DB) {
	app := bootstrap.App("../../../.env", "test")
	psqlDB := app.Postgres

	gin.SetMode(gin.TestMode)
	gin := gin.New()
	routerV1 := gin.Group("v1")
	routev1.Setup(app.Env, psqlDB, routerV1)

	return gin, psqlDB
}

func sampleUser() domain.User {
	return domain.User{
		Name:     "testUser",
		Email:    "testUser@gmail.com",
		Password: "testPassword",
	}
}

func signup(gin *gin.Engine, db *gorm.DB, user domain.User) (domain.SignupResponse, domain.User) {
	signupString := fmt.Sprintf(`{"name": "%s", "email": "%s", "password": "%s"}`, user.Name, user.Email, user.Password)
	signupReqBody := []byte(signupString)

	signupReq, _ := http.NewRequest(http.MethodPost, "/v1/signup", bytes.NewReader(signupReqBody))

	signupReq.Header.Add("Content-Type", "application/json")
	signupRec := httptest.NewRecorder()
	gin.ServeHTTP(signupRec, signupReq)

	signupResp := domain.SignupResponse{}
	json.NewDecoder(signupRec.Body).Decode(&signupResp)

	db.Model(&domain.User{}).Where("name = ?", user.Name).First(&user)

	return signupResp, user
}

func login(gin *gin.Engine, user domain.User) domain.LoginResponse {
	loginString := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, user.Email, user.Password)
	loginReqBody := []byte(loginString)

	loginReq, _ := http.NewRequest(http.MethodPost, "/v1/login", bytes.NewReader(loginReqBody))

	loginReq.Header.Add("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	gin.ServeHTTP(loginRec, loginReq)

	loginResp := domain.LoginResponse{}
	json.NewDecoder(loginRec.Body).Decode(&loginResp)

	return loginResp
}

func fetchUser(db *gorm.DB, userName string) (domain.User, error) {
	var user domain.User
	tx := db.Model(&domain.User{}).Where("name = ?", userName).First(&user)
	return user, tx.Error
}

func cleanupUser(db *gorm.DB, userName string) {
	var user domain.User
	db.Model(&domain.User{}).Where("name = ?", userName).First(&user)
	db.Unscoped().Delete(&user, "name = ?", userName)
}
