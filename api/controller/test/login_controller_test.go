package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	gin, db := setup()
	user := sampleUser()

	t.Run("success", func(t *testing.T) {
		signupReqBody := []byte(fmt.Sprintf(
			`{"name": "%s", "email": "%s", "password": "%s"}`,
			user.Name, user.Email, user.Password),
		)

		signupReq, err := http.NewRequest(http.MethodPost, "/v1/signup", bytes.NewReader(signupReqBody))

		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		signupReq.Header.Add("Content-Type", "application/json")
		signupRec := httptest.NewRecorder()
		gin.ServeHTTP(signupRec, signupReq)

		if signupRec.Code != http.StatusOK {
			t.Fatalf("Response returned with an unexpected status: %v\n", signupRec.Code)
		}

		loginReqBody := []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, user.Email, user.Password))

		loginReq, err := http.NewRequest(http.MethodPost, "/v1/login", bytes.NewReader(loginReqBody))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		loginReq.Header.Add("Content-Type", "application/json")
		loginRec := httptest.NewRecorder()
		gin.ServeHTTP(loginRec, loginReq)

		if loginRec.Code != http.StatusOK {
			t.Fatalf("Response returned with an unexpected status: %v\n", signupRec.Code)
		}
	})

	cleanupUser(db, user.Name)
}
