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
	ideas, _ := seedIdeas(db, user)

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

		ideaResp := []domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(&ideaResp)

		if len(ideaResp) != 2 {
			t.Fatalf("Unexpected response body length: %d\n", len(ideaResp))
		}
	})

	t.Run("byUserID", func(t *testing.T) {
		endpoint := fmt.Sprintf("/v1/users/%d/ideas", user.ID)
		ideaReq, err := http.NewRequest(http.MethodGet, endpoint, bytes.NewReader([]byte{}))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		ideaReq.Header.Add("Content-Type", "application/json")
		ideaRec := httptest.NewRecorder()

		gin.ServeHTTP(ideaRec, ideaReq)

		if ideaRec.Code != http.StatusOK {
			t.Fatalf("Unexpected response status code: %d\n", ideaRec.Code)
		}

		ideaResp := []domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(&ideaResp)

		if len(ideaResp) != 2 {
			t.Fatalf("Unexpected response body length: %d\n", len(ideaResp))
		}
	})

	t.Run("byID", func(t *testing.T) {
		endpoint := fmt.Sprintf("/v1/ideas/%d", ideas[0].ID)
		ideaReq, err := http.NewRequest(http.MethodGet, endpoint, bytes.NewReader([]byte{}))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		ideaReq.Header.Add("Content-Type", "application/json")
		ideaRec := httptest.NewRecorder()

		gin.ServeHTTP(ideaRec, ideaReq)

		if ideaRec.Code != http.StatusOK {
			t.Fatalf("Unexpected response status code: %d\n", ideaRec.Code)
		}

		ideaResp := domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(&ideaResp)

		if ideaResp.Content != ideas[0].Content {
			t.Fatalf("Unexpected response content:\n expected: %s\n obtained: %s\n", ideas[0].Content, ideaResp.Content)
		}
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
	gin, db := setup()
	authTokens := signup(gin, sampleUser())
	user, _ := fetchUser(db, sampleUser().Name)
	ideas, _ := seedIdeas(db, user)

	endpoint := fmt.Sprintf("/v1/ideas/%d", ideas[0].ID)
	ideaReq, err := http.NewRequest(http.MethodDelete, endpoint, bytes.NewReader([]byte{}))
	if err != nil {
		t.Fatalf("could not create request: %v\n", err)
	}

	ideaReq.Header.Add("Content-Type", "application/json")
	ideaReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
	ideaRec := httptest.NewRecorder()

	gin.ServeHTTP(ideaRec, ideaReq)

	if ideaRec.Code != http.StatusAccepted {
		t.Fatalf("Unexpected response status code: %d\n", ideaRec.Code)
	}

	remainingIdeas, _ := fetchIdeas(db)
	if len(ideas) == len(remainingIdeas) {
		t.Fatalf("Expected ideas slice length: %d, obtained: %d", len(ideas)-1, len(remainingIdeas))
	}

	cleanupIdeas(db)
	cleanupUser(db, user.Name)
}

func seedIdeas(db *gorm.DB, user domain.User) ([]domain.Idea, error) {
	firstIdea := domain.Idea{
		UserId:  user.ID,
		Content: "Some content that is suitable to a sample idea",
	}

	secondIdea := domain.Idea{
		UserId:  user.ID,
		Content: "Some other content which is still suitable",
	}

	tx := db.CreateInBatches([]domain.Idea{firstIdea, secondIdea}, 2)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var ideas []domain.Idea
	res := db.Find(&ideas)
	if res.Error != nil {
		return nil, res.Error
	}

	return ideas, nil
}

func fetchIdeas(db *gorm.DB) ([]domain.Idea, error) {
	var ideas []domain.Idea
	res := db.Find(&ideas)
	if res.Error != nil {
		return nil, res.Error
	}

	return ideas, nil
}

func cleanupIdeas(db *gorm.DB) {
	db.Exec("delete from ideas")
}
