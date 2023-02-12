package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	gin, db := setup()

	t.Run("success", func(t *testing.T) {
		signupBody := []byte(`{
			"name": "anchors",
			"email": "anchors@gmail.com",
			"password": "anchors"
		}`)

		signupBodyReader := bytes.NewReader(signupBody)
		signupReq, err := http.NewRequest(http.MethodPost, "/v1/signup", signupBodyReader)

		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		signupReq.Header.Add("Content-Type", "application/json")
		signupRec := httptest.NewRecorder()
		gin.ServeHTTP(signupRec, signupReq)

		if signupRec.Code != http.StatusOK {
			t.Fatalf("Response returned with an unexpected status: %v\n", signupRec.Code)
		}

		loginBody := []byte(`{
			"email": "anchors@gmail.com",
			"password": "anchors"
		}`)

		loginBodyReader := bytes.NewReader(loginBody)
		loginReq, err := http.NewRequest(http.MethodPost, "/v1/login", loginBodyReader)
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

	cleanupUser(db, "anchors")
}
