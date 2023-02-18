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
	gin, db := setup()
	user := sampleUser()
	signup(gin, user)
	user, _ = fetchUser(db, user.Name)
	seedIdeas(db, user)

	t.Run("all", func(t *testing.T) {
		ideaReq, err := http.NewRequest(http.MethodGet, "/v1/ideas", bytes.NewReader([]byte{}))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		ideaReq.Header.Add("Content-Type", "application/json")
		ideaRec := httptest.NewRecorder()

		gin.ServeHTTP(ideaRec, ideaReq)

		if ideaRec.Code != http.StatusOK {
			t.Fatalf("Response returned with an unexpected status code: %d\n", ideaRec.Code)
		}
	})

	t.Run("byID", func(t *testing.T) {
		t.Skip()
	})

	t.Run("byUserID", func(t *testing.T) {
		t.Skip()
	})

	cleanupIdeas(db)
	cleanupUser(db, user.Name)
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

func seedIdeas(db *gorm.DB, user domain.User) error {
	idea := domain.Idea{
		UserId:  user.ID,
		Content: "some content that is suitable to a sample idea",
	}

	tx := db.Create(&idea)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func cleanupIdeas(db *gorm.DB) {
	db.Exec("delete from ideas")
}
