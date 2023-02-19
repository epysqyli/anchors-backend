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

		ideas = fetchIdeas(db)
		if len(ideasResp) != len(ideas) {
			t.Fatalf("Expected ideas slice length: %d, obtained: %d", len(ideasResp), len(ideas))
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

	t.Cleanup(func() {
		cleanupDatabase(db)
		cleanupUser(db, user.Name)
	})
}

func TestCreateIdea(t *testing.T) {
	gin, db := setup()
	authTokens := signup(gin, sampleUser())

	// this test will eventually disappear as this behavior will not be permitted
	t.Run("noResources", func(t *testing.T) {
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

	t.Run("withResources", func(t *testing.T) {
		// arrange
		ideaReqBody := []byte(`{
			"content": "Idea with resource video",
			"videos": [
				{
					"url": "https://www.youtube.com/watch?v=8cX1aptP5Io&list=FL6zRqV5BoLaPshnUjI_oLPg&ab_channel=TheBitcoinLayer",
					"youtube_channel": "Some bitcoin channel"
				},
				{
					"url": "https://www.youtube.com/watch?v=MAeYCvyjQgE&ab_channel=JordanBPetersonClips",
					"youtube_channel": "Some bitcoin channel"
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

		// additional checks based on resources array resp and content
		ideaResp := domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(&ideaResp)

		if len(ideaResp.Videos) == 0 {
			t.Fatalf("No associated resources found")
		}
	})

	t.Run("withMultipleResourceTypes", func(t *testing.T) {
		// arrange
		ideaReqBody := []byte(`{
				"content": "Idea with video and blog resources",
				"videos": [
					{
						"url": "https://www.youtube.com/watch?v=8cX1aptP5Io&list=FL6zRqV5BoLaPshnUjI_oLPg&ab_channel=TheBitcoinLayer",
						"youtube_channel": "Some bitcoin channel"
					},
					{
						"url": "https://www.youtube.com/watch?v=MAeYCvyjQgE&ab_channel=JordanBPetersonClips",
						"youtube_channel": "Some bitcoin channel"
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
	})

	t.Run("preventDuplicateResource", func(t *testing.T) {
		// arrange
		url := "https://some-random-url.com"
		createBlog(db, url)
		blog := fetchBlogByUrl(db, url)
		blogsArray := fmt.Sprintf(`"blogs": [{"id": %d, "url": "%s", "category": "%s"}]`, blog.ID, url, blog.Category)

		ideaReqBody := []byte(fmt.Sprintf(
			`{"content": "Some random idea that I'd like to publish", %s}`, blogsArray))

		anotherIdeaReqBody := []byte(fmt.Sprintf(
			`{"content": "Some random idea that I'd like to publish", %s}`, blogsArray))

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

		previousBlogsCount := len(fetchBlogs(db))

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

		currentBlogsCount := len(fetchBlogs(db))

		if previousBlogsCount != currentBlogsCount {
			t.Fatalf("Number of blogs increased from %d to %d", previousBlogsCount, currentBlogsCount)
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

	remainingIdeas := fetchIdeas(db)
	if len(ideas) == len(remainingIdeas) {
		t.Fatalf("Expected ideas slice length: %d, obtained: %d", len(ideas)-1, len(remainingIdeas))
	}

	t.Cleanup(func() {
		cleanupDatabase(db)
		cleanupUser(db, user.Name)
	})
}

func seedIdeas(db *gorm.DB, user domain.User) []domain.Idea {
	firstIdea := domain.Idea{
		UserID:  user.ID,
		Content: "Some content that is suitable to a sample idea",
	}

	secondIdea := domain.Idea{
		UserID:  user.ID,
		Content: "Some other content which is still suitable",
	}

	db.CreateInBatches([]domain.Idea{firstIdea, secondIdea}, 2)

	ideas := fetchIdeas(db)
	return ideas
}

func fetchIdeas(db *gorm.DB) []domain.Idea {
	var ideas []domain.Idea
	db.Find(&ideas)

	return ideas
}

func createBlog(db *gorm.DB, url string) {
	blog := domain.Blog{
		Url:      url,
		Category: "software development",
	}

	db.Create(&blog)
}

func fetchBlogByUrl(db *gorm.DB, url string) domain.Blog {
	var blog domain.Blog
	db.Model(&domain.Blog{}).Where("url = ?", url).First(&blog)
	return blog
}

// make a generic fetch func?
func fetchBlogs(db *gorm.DB) []domain.Blog {
	var blogs []domain.Blog
	db.Find(&blogs)

	return blogs
}

func cleanupDatabase(db *gorm.DB) {
	db.Exec("delete from ideas_videos")
	db.Exec("delete from videos")
	db.Exec("delete from blogs_ideas")
	db.Exec("delete from blogs")
	db.Exec("delete from ideas")
}
