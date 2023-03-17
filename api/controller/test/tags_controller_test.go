package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
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

	t.Run("FetchAllTags", func(t *testing.T) {
		t.Skip()
	})

	t.Run("FetchByID", func(t *testing.T) {
		t.Skip()
	})

	t.Run("FetchByName", func(t *testing.T) {
		t.Skip()
	})

	t.Run("DeletebyID", func(t *testing.T) {
		t.Skip()
	})

	t.Cleanup(func() {
		cleanupDatabase(db)
		cleanupUser(db, user.Name)
	})
}
