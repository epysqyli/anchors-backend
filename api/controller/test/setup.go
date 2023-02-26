package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	routev1 "github.com/epysqyli/anchors-backend/api/route/v1"
	"github.com/epysqyli/anchors-backend/bootstrap"
	"github.com/epysqyli/anchors-backend/domain"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TestMain for controller test package?
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
