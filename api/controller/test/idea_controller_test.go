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
	_, user := signup(gin, db, sampleUser())
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

		ideas := fetchResources(db, []domain.Idea{})
		if len(ideasResp) != len(ideas) {
			t.Fatalf("Expected ideas slice length: %d, obtained: %d", len(ideasResp), len(ideas))
		}

		checkIdeaAssociations(t, &ideasResp[1])
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

		if len(ideaResp) != len(ideas) {
			t.Fatalf("Unexpected response body length: %d\n", len(ideaResp))
		}

		checkIdeaAssociations(t, &ideaResp[1])
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

		checkIdeaAssociations(t, &ideaResp)
	})

	t.Cleanup(func() {
		cleanupDatabase(db)
		cleanupUser(db, user.Name)
	})
}

func TestCreateIdeas(t *testing.T) {
	gin, db := setup()
	authTokens, user := signup(gin, db, sampleUser())

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

	// https://www.youtube.com/watch?v=p8u_k2LIZyo
	t.Run("withVideos", func(t *testing.T) {
		ideaReqBody := []byte(`{
			"content": "Idea with video and blog resources",
			"videos": [
				{"url": "https://www.youtube.com/watch?v=p8u_k2LIZyo"},
				{"url": "https://www.youtube.com/watch?v=tKbV6BpH-C8&t=167s&ab_channel=CodeAesthetic"},
				{"url": "https://www.randomvideos.com/videos/12444"},
				{"url": "https://youtu.be/7_cXxEbR_pA"}
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

		idea := &domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(idea)

		if idea.Videos[0].Identifier != "p8u_k2LIZyo" {
			t.Fatalf("Wrong identifier\nExpected: %s\nReturned: %s", "p8u_k2LIZyo", idea.Videos[0].Identifier)
		}

		if idea.Videos[1].Identifier != "tKbV6BpH-C8" {
			t.Fatalf("Wrong identifier\nExpected: %s\nReturned: %s", "tKbV6BpH-C8", idea.Videos[1].Identifier)
		}

		nonYtIdentifier := "https://www.randomvideos.com/videos/12444"
		if idea.Videos[2].Identifier != nonYtIdentifier {
			t.Fatalf("Wrong identifier\nExpected: %s\nReturned: %s", nonYtIdentifier, idea.Videos[2].Identifier)
		}

		if idea.Videos[3].Identifier != "7_cXxEbR_pA" {
			t.Fatalf("Wrong identifier\nExpected: %s\nReturned: %s", "7_cXxEbR_pA", idea.Videos[3].Identifier)
		}
	})

	t.Run("withMultipleResourceTypes", func(t *testing.T) {
		// arrange
		ideaReqBody := []byte(`{
				"content": "Idea with video and blog resources",
				"videos": [
					{
						"url": "https://www.youtube.com/watch?v=8cX1aptP5Io&list=FL6zRqV5BoLaPshnUjI_oLPg&ab_channel=TheBitcoinLayer",
						"timestamp": 124
					},
					{
						"url": "https://www.youtube.com/watch?v=MAeYCvyjQgE&ab_channel=JordanBPetersonClips",
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

		// uncomment once youtube unique identifier issue is solved
		// if ideaResp.Videos[0].YoutubeChannel == "" || ideaResp.Videos[1].YoutubeChannel == "" {
		// 	t.Fatal("Youtube Channel for videos not assigned")
		// }

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
		ideaReqBody := []byte(fmt.Sprintf((`{"content": "New idea with two anchor ideas",
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

	t.Run("withBookAndChapter", func(t *testing.T) {
		bookResources := `[
			{
				"url": "https://openlibrary.org/works/OL20984004W",
				"open_library_key": "OL20984004W",
				"title": "The Bitcoin Standard",
				"year": 2018,
				"number_of_pages": 304,
				"open_library_id": 10320866,
				"language": "eng",
				"authors": [
					{"open_library_key": "OL7945937A", "full_name": "James Fouhey"},
					{"open_library_key": "OL8027052A", "full_name": "Saifedean Ammous"}
				],
				"chapter": "2 - the greatest chapter of all time"
			},
			{
				"url": "https://openlibrary.org/works/OL20984100F",
				"open_library_key": "OL20984100F",
				"title": "The Whatever Book Title",
				"year": 2000,
				"number_of_pages": 200,
				"open_library_id": 10320100,
				"language": "eng",
				"authors": [{"open_library_key": "OL1005931M", "full_name": "Best Writer"}]
			}
		]`

		ideaReqBody := []byte(fmt.Sprintf((`{"content": "Idea based on a book", "books": %s}`), bookResources))

		ideaReq, err := http.NewRequest(http.MethodPost, "/v1/ideas", bytes.NewReader(ideaReqBody))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		ideaReq.Header.Add("Content-Type", "application/json")
		ideaReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		ideaRec := httptest.NewRecorder()
		gin.ServeHTTP(ideaRec, ideaReq)

		if ideaRec.Code != http.StatusCreated {
			t.Fatalf("Response returned with an unexpected status code: %d\n", ideaRec.Code)
		}

		ideaResp := &domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(ideaResp)

		firstBookIdeaRel := &domain.BooksIdeas{
			IdeaID: ideaResp.ID,
			BookID: ideaResp.Books[0].ID,
		}
		db.Find(firstBookIdeaRel)

		if firstBookIdeaRel.Chapter == "" {
			t.Fatal("Chapter field from the book idea relation is empty")
		}

		secondBookIdeaRel := &domain.BooksIdeas{
			IdeaID: ideaResp.ID,
			BookID: ideaResp.Books[1].ID,
		}
		db.Find(secondBookIdeaRel)

		if secondBookIdeaRel.Chapter != "" {
			t.Fatal("Chapter field from the book idea relation should be empty")
		}
	})

	t.Run("twoIdeasSameBook", func(t *testing.T) {
		bookResource := `{
			"url": "https://openlibrary.org/works/OL00000000A",
			"open_library_key": "OL00000000A",
			"title": "Basic book title",
			"year": 2000,
			"number_of_pages": 100,
			"open_library_id": 10101010,
			"language": "eng",
			"authors": [
				{"open_library_key": "OL8051252D", "full_name": "Basic Author"}
			]
		}`

		firstIdea := fmt.Sprintf((`{"content": "Some idea - book one", "books": [%s]}`), bookResource)
		firstIdeaReq, _ := http.NewRequest(http.MethodPost, "/v1/ideas", bytes.NewReader([]byte(firstIdea)))
		firstIdeaReq.Header.Add("Content-Type", "application/json")
		firstIdeaReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		firstIdeaRec := httptest.NewRecorder()
		gin.ServeHTTP(firstIdeaRec, firstIdeaReq)

		ideaResp := &domain.Idea{}
		json.NewDecoder(firstIdeaRec.Body).Decode(ideaResp)
		bookID := ideaResp.Books[0].ID

		secondIdea := fmt.Sprintf((`{"content": "Another idea - book one", "books": [%s]}`), bookResource)
		secondIdeaReq, _ := http.NewRequest(http.MethodPost, "/v1/ideas", bytes.NewReader([]byte(secondIdea)))
		secondIdeaReq.Header.Add("Content-Type", "application/json")
		secondIdeaReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		secondIdeaRec := httptest.NewRecorder()
		gin.ServeHTTP(secondIdeaRec, secondIdeaReq)

		// check that many to many has two entries for the book
		bookIdeasRels := []domain.BooksIdeas{}
		db.Where([]domain.BooksIdeas{{BookID: bookID}}).Find(&bookIdeasRels)

		if len(bookIdeasRels) != 2 {
			t.Fatalf("books_ideas entries for bookID %d\nexpected: %d, got: %d", bookID, 2, len(bookIdeasRels))
		}
	})

	t.Cleanup(func() {
		cleanupDatabase(db)
		cleanupUser(db, sampleUser().Name)
	})
}

func TestDeleteIdea(t *testing.T) {
	gin, db := setup()
	authTokens, user := signup(gin, db, sampleUser())
	ideas := seedIdeas(db, user)

	endpoint := fmt.Sprintf("/v1/ideas/%d", ideas[1].ID)
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

	// test fetch idea -> should be empty
	var idea domain.Idea
	db.First(&idea, ideas[1].ID)

	if idea.ID != 0 {
		t.Fatalf("Idea with ID: %d should not exist", ideas[1].ID)
	}

	t.Cleanup(func() {
		cleanupDatabase(db)
		cleanupUser(db, user.Name)
	})
}

func seedIdeas(db *gorm.DB, user domain.User) []domain.Idea {
	emptyIdea := &domain.Idea{
		UserID:  user.ID,
		Content: "Some content that is suitable to a sample idea",
	}

	db.Create(emptyIdea)

	book := domain.Book{
		Url:            "https://openlibrary.org/works/OL02600010B",
		OpenLibraryKey: "OL02600010B",
		Title:          "Book from seeds",
		Year:           1999,
		NumberOfPages:  100,
		OpenLibraryID:  12341234,
		Language:       "eng",
		Authors: []domain.Author{
			{OpenLibraryKey: "OL8054255E", FullName: "Some Writer"},
		},
	}

	fullIdea := &domain.Idea{
		UserID:  user.ID,
		Content: "Content for an idea anchored upon a blog",
		Blogs:   []domain.Blog{{Url: "https://some-blog.com", Category: "science"}},
		Videos:  []domain.Video{{Url: "https://some-youtube-video.com", YoutubeChannel: "cool-channel"}},
		Anchors: []*domain.Idea{emptyIdea},
		Books:   []domain.Book{book},
	}

	db.Create(fullIdea)

	return []domain.Idea{*emptyIdea, *fullIdea}
}

func fetchResources[M any](db *gorm.DB, resources []M) []M {
	db.Find(&resources)
	return resources
}

func fetchResourceByUrl[M any](db *gorm.DB, resource *M, url string) *M {
	db.Model(resource).Where("url = ?", url).First(resource)
	return resource
}

func cleanupDatabase(db *gorm.DB) {
	db.Exec("delete from anchors_ideas")
	db.Exec("delete from ideas_videos")
	db.Exec("delete from blogs_ideas")
	db.Exec("delete from books_ideas")
	db.Exec("delete from authors_books")
	db.Exec("delete from authors")
	db.Exec("delete from books")
	db.Exec("delete from videos")
	db.Exec("delete from blogs")
	db.Exec("delete from ideas")
}

func checkIdeaAssociations(t *testing.T, idea *domain.Idea) {
	if len(idea.Blogs) == 0 {
		t.Fatalf("Blogs missing, expected: %d, obtained: %d", 1, len(idea.Blogs))
	}

	if len(idea.Videos) == 0 {
		t.Fatalf("Videos missing, expected: %d, obtained: %d", 1, len(idea.Videos))
	}

	if len(idea.Anchors) == 0 {
		t.Fatalf("Anchor ideas missing, expected: %d, obtained: %d", 1, len(idea.Anchors))
	}

	if len(idea.Books) == 0 {
		t.Fatalf("Books missing, expected: %d, obtained: %d", 1, len(idea.Books))
	}

	if len(idea.Books[0].Authors) == 0 {
		t.Fatalf("Author missing on books, expected: %d, obtained: %d", 1, len(idea.Books[0].Authors))
	}
}
