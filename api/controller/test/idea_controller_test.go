package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchIdeas(t *testing.T) {
	// gin, db := setup()
	// 1st: create new ideas associated to the user
	// 2nd: create ideas with resources associated to them and test for the expected structure
	// add a Book and a youtube video as resources in the second test

	t.Run("byID", func(t *testing.T) {
		t.Skip()
	})

	t.Run("byUserID", func(t *testing.T) {
		t.Skip()
	})

	t.Run("allIdeas", func(t *testing.T) {
		t.Skip()
	})

	// teardown
}

func TestCreateIdea(t *testing.T) {
	gin, db := setup()
	authTokens := signup(gin, sampleUser())

	t.Run("basicIdea", func(t *testing.T) {
		// arrange
		ideaReqBody := []byte(`{"content": "this is a test idea"}`)
		ideaReq, err := http.NewRequest(http.MethodPost, "/v1/ideas", bytes.NewReader(ideaReqBody))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		ideaReq.Header.Add("Content-Type", "application/json")
		ideaReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		ideaRec := httptest.NewRecorder()

		// act
		gin.ServeHTTP(ideaRec, ideaReq)

		// assert
		if ideaRec.Code != http.StatusCreated {
			t.Fatalf("Response returned with an unexpected status code: %d\n", ideaRec.Code)
		}
	})

	t.Run("ideaWithOneResource", func(t *testing.T) {
		t.Skip()
	})

	t.Run("ideaWithMultipleResource", func(t *testing.T) {
		t.Skip()
	})

	// cleanup idea and other resources from db
	cleanupUser(db, sampleUser().Name)
}

func TestDeleteIdeaByID(t *testing.T) {
	// logged in user is needed
	t.Skip()
}
