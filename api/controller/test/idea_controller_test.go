package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/epysqyli/anchors-backend/domain"
	"gorm.io/gorm"
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
		content := "this is a test idea"
		ideaReqBody := []byte(fmt.Sprintf(`{"content": "%s"}`, content))
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

		ideaResp := domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(&ideaResp)

		if ideaResp.Content != content {
			t.Fatalf("Response returned with an unexpected content: \texpected: %s\n\tobtained: %s\n",
				content, ideaResp.Content)
		}
	})

	t.Run("ideaWithOneResource", func(t *testing.T) {
		t.Skip()
	})

	t.Run("ideaWithMultipleResource", func(t *testing.T) {
		t.Skip()
	})

	cleanupIdeas(db)
	cleanupUser(db, sampleUser().Name)
}

func TestDeleteIdeaByID(t *testing.T) {
	// logged in user is needed
	t.Skip()
}

func cleanupIdeas(db *gorm.DB) {
	db.Exec("delete from ideas")
}
