package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignup(t *testing.T) {
	gin, db := setup()
	user := sampleUser()

	t.Run("success", func(t *testing.T) {
		signupReqBody := []byte(fmt.Sprintf(
			`{"name": "%s", "email": "%s", "password": "%s"}`,
			user.Name, user.Email, user.Password),
		)

		req, err := http.NewRequest(http.MethodPost, "/v1/signup", bytes.NewReader(signupReqBody))

		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		req.Header.Add("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		gin.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("Response returned with an unexpected status: %v\n", rec.Code)
		}

	})

	cleanupUser(db, user.Name)
}
