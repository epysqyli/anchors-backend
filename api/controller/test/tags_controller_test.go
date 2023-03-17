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

func TestTags(t *testing.T) {
	gin, db := setup()
	authTokens, user := signup(gin, db, sampleUser())

	t.Run("Create", func(t *testing.T) {
		tagReqBody := []byte(`{"name": "economics"}`)
		tagReq, err := http.NewRequest(http.MethodPost, "/v1/tags", bytes.NewReader(tagReqBody))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		tagReq.Header.Add("Content-Type", "application/json")
		tagReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		tagRec := httptest.NewRecorder()

		gin.ServeHTTP(tagRec, tagReq)
		assertEqual(http.StatusCreated, tagRec.Code, t, "Tag should have been created")
	})

	t.Run("FetchAll", func(t *testing.T) {
		createTag(db, "architecture")
		tagReq, err := http.NewRequest(http.MethodGet, "/v1/tags", bytes.NewReader([]byte{}))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		tagReq.Header.Add("Content-Type", "application/json")
		tagReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		tagRec := httptest.NewRecorder()

		gin.ServeHTTP(tagRec, tagReq)
		assertEqual(http.StatusOK, tagRec.Code, t)

		tags := []domain.Tag{}
		json.NewDecoder(tagRec.Body).Decode(&tags)
		assertUnequal(0, len(tags), t, "Tags should have been fetched")
	})

	t.Run("FetchByID", func(t *testing.T) {
		someTag := createTag(db, "economics")
		tagReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/tags/%d", someTag.ID), bytes.NewReader([]byte{}))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		tagReq.Header.Add("Content-Type", "application/json")
		tagReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		tagRec := httptest.NewRecorder()

		gin.ServeHTTP(tagRec, tagReq)
		assertEqual(http.StatusOK, tagRec.Code, t)

		tag := domain.Tag{}
		json.NewDecoder(tagRec.Body).Decode(&tag)
		assertUnequal(0, tag.ID, t, "Tag ID should be assigned")
	})

	t.Run("FetchByName", func(t *testing.T) {
		someTag := createTag(db, "finance")
		tagReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/tags/%s", someTag.Name), bytes.NewReader([]byte{}))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		tagReq.Header.Add("Content-Type", "application/json")
		tagReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		tagRec := httptest.NewRecorder()

		gin.ServeHTTP(tagRec, tagReq)
		assertEqual(http.StatusOK, tagRec.Code, t)

		tag := domain.Tag{}
		json.NewDecoder(tagRec.Body).Decode(&tag)
		assertUnequal(0, tag.Name, t, "Tag Name should be assigned")
	})

	t.Run("DeletebyID", func(t *testing.T) {
		someTag := createTag(db, "del")
		tagReq, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/tags/%d", someTag.ID), bytes.NewReader([]byte{}))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		tagReq.Header.Add("Content-Type", "application/json")
		tagReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		tagRec := httptest.NewRecorder()

		gin.ServeHTTP(tagRec, tagReq)
		assertEqual(http.StatusAccepted, tagRec.Code, t)

		tag := domain.Tag{}
		db.Where(&domain.Tag{Name: "del"}).First(&tag)
		assertEqual(uint(0), tag.ID, t, "There should be no 'del' tag")
	})

	t.Cleanup(func() {
		cleanupDatabase(db)
		cleanupUser(db, user.Name)
	})
}

func createTag(db *gorm.DB, name string) domain.Tag {
	tag := domain.Tag{Name: name}
	db.Create(&tag)

	return tag
}

func fetchTag(db *gorm.DB, name string) *domain.Tag {
	tag := &domain.Tag{}
	db.Where(&domain.Tag{Name: name}).First(tag)

	return tag
}
