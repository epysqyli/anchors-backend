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
			t.Fatalf("\n\n  expected: %v\n  obtained: %v\n\n", exp, got)
		} else {
			t.Fatalf("\n\n%s\n  expected: %v\n  obtained: %v\n\n", info[0], exp, got)
		}
	}
}

func assertUnequal(exp any, got any, t *testing.T, info ...string) {
	if exp == got {
		if len(info) == 0 {
			t.Fatalf("\n\n  %c expected: %v\n  obtained: %v\n\n", rune(172), exp, got)
		} else {
			t.Fatalf("\n\n%s\n  %c expected: %v\n  obtained: %v\n\n", info[0], rune(172), exp, got)
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

func cleanupDatabase(db *gorm.DB) {
	db.Exec("delete from musical_artists_songs")
	db.Exec("delete from ideas_tags")
	db.Exec("delete from generics_ideas")
	db.Exec("delete from articles_ideas")
	db.Exec("delete from ideas_wikis")
	db.Exec("delete from anchors_ideas")
	db.Exec("delete from ideas_songs")
	db.Exec("delete from ideas_videos")
	db.Exec("delete from blogs_ideas")
	db.Exec("delete from books_ideas")
	db.Exec("delete from authors_books")
	db.Exec("delete from ideas_movies")
	db.Exec("delete from cinematic_genres_movies")
	db.Exec("delete from cinematic_genres")
	db.Exec("delete from musical_artists")
	db.Exec("delete from tags")
	db.Exec("delete from generics")
	db.Exec("delete from articles")
	db.Exec("delete from songs")
	db.Exec("delete from wikis")
	db.Exec("delete from musical_albums")
	db.Exec("delete from movies")
	db.Exec("delete from authors")
	db.Exec("delete from books")
	db.Exec("delete from videos")
	db.Exec("delete from blogs")
	db.Exec("delete from ideas")
}
