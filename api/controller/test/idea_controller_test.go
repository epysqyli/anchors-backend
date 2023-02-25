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
	ideas := seedIdeas(db, user)

	// check for associated resources
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

		ideasResp := []domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(&ideasResp)

		ideas := fetchResources(db, []domain.Idea{})
		if len(ideasResp) != len(ideas) {
			t.Fatalf("Expected ideas slice length: %d, obtained: %d", len(ideasResp), len(ideas))
		}
	})

	// check for associated resources
	t.Run("byUserID", func(t *testing.T) {
		t.Skip()
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
		endpoint := fmt.Sprintf("/v1/ideas/%d", ideas[1].ID)
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

		if ideaResp.Content != ideas[1].Content {
			t.Fatalf("Unexpected response content:\n expected: %s\n obtained: %s\n", ideas[0].Content, ideaResp.Content)
		}

		if len(ideaResp.Blogs) == 0 {
			t.Fatalf("Blogs missing, expected: %d, obtained: %d", 1, len(ideaResp.Blogs))
		}

		if len(ideaResp.Videos) == 0 {
			t.Fatalf("Videos missing, expected: %d, obtained: %d", 1, len(ideaResp.Videos))
		}

		if len(ideaResp.Anchors) == 0 {
			t.Fatalf("Anchor ideas missing, expected: %d, obtained: %d", 1, len(ideaResp.Anchors))
		}
	})

	t.Cleanup(func() {
		cleanupDatabase(db)
		cleanupUser(db, user.Name)
	})
}

func TestCreateIdeas(t *testing.T) {
	gin, db := setup()
	authTokens := signup(gin, sampleUser())
	user, _ := fetchUser(db, sampleUser().Name)

	t.Run("withNoResources", func(t *testing.T) {
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
		if ideaRec.Code != http.StatusBadRequest {
			t.Fatalf("Response returned with an unexpected status code: %d\n", ideaRec.Code)
		}
	})

	t.Run("withMultipleResourceTypes", func(t *testing.T) {
		// arrange
		ideaReqBody := []byte(`{
				"content": "Idea with video and blog resources",
				"videos": [
					{
						"url": "https://www.youtube.com/watch?v=8cX1aptP5Io&list=FL6zRqV5BoLaPshnUjI_oLPg&ab_channel=TheBitcoinLayer",
						"youtube_channel": "Some bitcoin channel",
						"timestamp": 124
					},
					{
						"url": "https://www.youtube.com/watch?v=MAeYCvyjQgE&ab_channel=JordanBPetersonClips",
						"youtube_channel": "Some bitcoin channel",
						"timestamp": 99
					}
				],
				"blogs": [
					{
						"url": "https://mtlynch.io/solo-developer-year-5/",
						"category": "software development"
					},
					{
						"url": "https://matt-rickard.com/ask-dumb-questions",
						"category": "software development"
					}
				]
			}`)

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

		if len(ideaResp.Videos) == 0 {
			t.Fatalf("No videos found")
		}

		if len(ideaResp.Blogs) == 0 {
			t.Fatalf("No blogs found")
		}

		firstIdeaVideoRelation := domain.IdeasVideos{
			IdeaID:  ideaResp.ID,
			VideoID: ideaResp.Videos[0].ID,
		}
		db.Find(&firstIdeaVideoRelation)

		if firstIdeaVideoRelation.Timestamp != 124 {
			t.Fatalf("Timestamp not correctly assigned\n\texpected: %d\n\tgot: %d",
				124, firstIdeaVideoRelation.Timestamp)
		}

		secondIdeaVideoRelation := domain.IdeasVideos{
			IdeaID:  ideaResp.ID,
			VideoID: ideaResp.Videos[1].ID,
		}
		db.Find(&secondIdeaVideoRelation)

		if secondIdeaVideoRelation.Timestamp != 99 {
			t.Fatalf("Timestamp not correctly assigned\n\texpected: %d\n\tgot: %d",
				99, secondIdeaVideoRelation.Timestamp)
		}

		blogsIdeasRelations := []domain.BlogsIdeas{}
		db.Find(&blogsIdeasRelations)

		if len(blogsIdeasRelations) != 2 {
			t.Fatalf("Wrong number of blogs ideas relations\n\texpected: %d\n\tgot: %d", 2, len(blogsIdeasRelations))
		}
	})

	t.Run("withExistingResource", func(t *testing.T) {
		// arrange
		db.Create(&domain.Video{Url: "https://some-random-url.com", YoutubeChannel: "some-channel"})
		video := fetchResourceByUrl(db, &domain.Video{}, "https://some-random-url.com")
		videoArray := fmt.Sprintf(`"videos": [{"id": %d, "url": "%s", "youtube_channel": "%s"}]`,
			video.ID, video.Url, video.YoutubeChannel)

		ideaReqBody := []byte(fmt.Sprintf(
			`{"content": "Some random idea that I'd like to publish", %s}`, videoArray))

		anotherIdeaReqBody := []byte(fmt.Sprintf(
			`{"content": "Some random idea that I'd like to publish", %s}`, videoArray))

		ideaReq, err := http.NewRequest(http.MethodPost, "/v1/ideas", bytes.NewReader(ideaReqBody))
		anotherIdeaReq, err := http.NewRequest(http.MethodPost, "/v1/ideas", bytes.NewReader(anotherIdeaReqBody))

		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		ideaReq.Header.Add("Content-Type", "application/json")
		ideaReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		ideaRec := httptest.NewRecorder()

		anotherIdeaReq.Header.Add("Content-Type", "application/json")
		anotherIdeaReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		anotherIdeaRec := httptest.NewRecorder()

		previousVideosCount := len(fetchResources(db, []domain.Video{}))

		// act
		gin.ServeHTTP(ideaRec, ideaReq)
		gin.ServeHTTP(anotherIdeaRec, anotherIdeaReq)

		// assert
		if ideaRec.Code != http.StatusCreated {
			t.Fatalf("Response returned with an unexpected status code: %d\n", ideaRec.Code)
		}

		if anotherIdeaRec.Code != http.StatusCreated {
			t.Fatalf("Response returned with an unexpected status code: %d\n", ideaRec.Code)
		}

		currentVideosCount := len(fetchResources(db, []domain.Video{}))

		if previousVideosCount != currentVideosCount {
			t.Fatalf("Number of blogs increased from %d to %d", previousVideosCount, currentVideosCount)
		}
	})

	t.Run("withAnchorIdeas", func(t *testing.T) {
		// arrange
		ideas := seedIdeas(db, user)
		ideaReqBody := []byte(fmt.Sprintf((`{"content": "Idea with video and blog resources",
			"anchors": [{"id": %d}, {"id": %d}]}`), ideas[0].ID, ideas[1].ID))

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

		ideaResp := &domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(ideaResp)

		if len(ideaResp.Anchors) != len(ideas) {
			t.Fatalf("Expected %d anchor ideas, obtained %d\n", len(ideas), len(ideaResp.Anchors))
		}
	})

	t.Cleanup(func() {
		cleanupDatabase(db)
		cleanupUser(db, sampleUser().Name)
	})
}

func TestDeleteIdeaByID(t *testing.T) {
	gin, db := setup()
	authTokens := signup(gin, sampleUser())
	user, _ := fetchUser(db, sampleUser().Name)
	ideas := seedIdeas(db, user)

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

	remainingIdeas := fetchResources(db, []domain.Idea{})
	if len(ideas) == len(remainingIdeas) {
		t.Fatalf("Expected ideas slice length: %d, obtained: %d", len(ideas)-1, len(remainingIdeas))
	}

	t.Cleanup(func() {
		cleanupDatabase(db)
		cleanupUser(db, user.Name)
	})
}

func seedIdeas(db *gorm.DB, user domain.User) []domain.Idea {
	emptyIdea := domain.Idea{
		UserID:  user.ID,
		Content: "Some content that is suitable to a sample idea",
	}

	fullIdea := domain.Idea{
		UserID:  user.ID,
		Content: "Content for an idea anchored upon a blog",
		Blogs:   []domain.Blog{{Url: "https://some-blog.com", Category: "science"}},
		Videos:  []domain.Video{{Url: "https://some-youtube-video.com", YoutubeChannel: "cool-channel"}},
		Anchors: []domain.Idea{emptyIdea},
	}

	db.Create(&emptyIdea)
	db.Create(&fullIdea)

	return []domain.Idea{emptyIdea, fullIdea}
}

func fetchResources[M any](db *gorm.DB, resources []M) []M {
	db.Find(&resources)
	return resources
}

func fetchResourceByUrl[M any](db *gorm.DB, resource *M, url string) *M {
	db.Model(resource).Where("url = ?", url).First(&resource)
	return resource
}

func cleanupDatabase(db *gorm.DB) {
	db.Exec("delete from ideas_anchors")
	db.Exec("delete from ideas_videos")
	db.Exec("delete from blogs_ideas")
	db.Exec("delete from videos")
	db.Exec("delete from blogs")
	db.Exec("delete from ideas")
}
