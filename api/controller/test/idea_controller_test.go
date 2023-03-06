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

	t.Run("All", func(t *testing.T) {
		ideaReq, err := http.NewRequest(http.MethodGet, "/v1/ideas", bytes.NewReader([]byte{}))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		ideaReq.Header.Add("Content-Type", "application/json")
		ideaRec := httptest.NewRecorder()

		gin.ServeHTTP(ideaRec, ideaReq)
		assertEqual(ideaRec.Code, http.StatusOK, t, "Idea should have been fetched")

		ideasResp := []domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(&ideasResp)

		ideas := fetchResources(db, []domain.Idea{})
		assertEqual(len(ideas), len(ideasResp), t, "Wrong number of ideas fetched")

		checkIdeaAssociations(t, &ideasResp[1])
	})

	t.Run("UserID", func(t *testing.T) {
		endpoint := fmt.Sprintf("/v1/users/%d/ideas", user.ID)
		ideaReq, err := http.NewRequest(http.MethodGet, endpoint, bytes.NewReader([]byte{}))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		ideaReq.Header.Add("Content-Type", "application/json")
		ideaRec := httptest.NewRecorder()

		gin.ServeHTTP(ideaRec, ideaReq)
		assertEqual(ideaRec.Code, http.StatusOK, t, "Idea should have been fetched")

		ideaResp := []domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(&ideaResp)
		assertEqual(len(ideas), len(ideaResp), t, "Wrong number of ideas fetched")

		checkIdeaAssociations(t, &ideaResp[1])
	})

	t.Run("ID", func(t *testing.T) {
		endpoint := fmt.Sprintf("/v1/ideas/%d", ideas[1].ID)
		ideaReq, err := http.NewRequest(http.MethodGet, endpoint, bytes.NewReader([]byte{}))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		ideaReq.Header.Add("Content-Type", "application/json")
		ideaRec := httptest.NewRecorder()

		gin.ServeHTTP(ideaRec, ideaReq)
		assertEqual(ideaRec.Code, http.StatusOK, t, "Idea should have been fetched")

		ideaResp := domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(&ideaResp)
		assertEqual(ideas[1].Content, ideaResp.Content, t, "Unexpected content")

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

	t.Run("Empty", func(t *testing.T) {
		ideaReqBody := []byte(`{"content": "this is a test idea"}`)
		ideaReq, err := http.NewRequest(http.MethodPost, "/v1/ideas", bytes.NewReader(ideaReqBody))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		ideaReq.Header.Add("Content-Type", "application/json")
		ideaReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		ideaRec := httptest.NewRecorder()

		gin.ServeHTTP(ideaRec, ideaReq)
		assertEqual(http.StatusBadRequest, ideaRec.Code, t, "Idea should not have been created")
	})

	t.Run("Videos", func(t *testing.T) {
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

		gin.ServeHTTP(ideaRec, ideaReq)
		assertEqual(ideaRec.Code, http.StatusCreated, t, "Idea should have been created")

		idea := &domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(idea)

		assertEqual("p8u_k2LIZyo", idea.Videos[0].Identifier, t, "Wrong identifier")
		assertEqual("tKbV6BpH-C8", idea.Videos[1].Identifier, t, "Wrong identifier")
		assertEqual("https://www.randomvideos.com/videos/12444", idea.Videos[2].Identifier, t, "Wrong identifier")
		assertEqual("7_cXxEbR_pA", idea.Videos[3].Identifier, t, "Wrong identifier")
	})

	t.Run("VideosWithTimestamp", func(t *testing.T) {
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
				]
			}`)

		ideaReq, err := http.NewRequest(http.MethodPost, "/v1/ideas", bytes.NewReader(ideaReqBody))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		ideaReq.Header.Add("Content-Type", "application/json")
		ideaReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		ideaRec := httptest.NewRecorder()

		gin.ServeHTTP(ideaRec, ideaReq)
		assertEqual(http.StatusCreated, ideaRec.Code, t, "Idea should have been created")

		ideaResp := domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(&ideaResp)

		assertUnequal(0, ideaResp.Videos, t, "No videos found")

		firstIdeaVideoRelation := domain.IdeasVideos{
			IdeaID:  ideaResp.ID,
			VideoID: ideaResp.Videos[0].ID,
		}
		db.Find(&firstIdeaVideoRelation)

		secondIdeaVideoRelation := domain.IdeasVideos{
			IdeaID:  ideaResp.ID,
			VideoID: ideaResp.Videos[1].ID,
		}
		db.Find(&secondIdeaVideoRelation)

		assertEqual(int16(124), firstIdeaVideoRelation.Timestamp, t, "Timestamp not correcly assigned")
		assertEqual(int16(99), secondIdeaVideoRelation.Timestamp, t, "Timestamp not correcly assigned")
	})

	t.Run("MultipleTypes", func(t *testing.T) {
		ideaReqBody := []byte(`{
			"content": "Idea with video and blog resources",
			"videos": [
				{"url": "https://www.youtube.com/watch?v=f2a_k2LIZyo"},
				{"url": "https://www.randomvideos.com/videos/11111"}
			],
			"blogs": [
				{"url": "https://cool-blog.com", "category": "low-level-programming"}
			]
		}`)

		ideaReq, err := http.NewRequest(http.MethodPost, "/v1/ideas", bytes.NewReader(ideaReqBody))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		ideaReq.Header.Add("Content-Type", "application/json")
		ideaReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		ideaRec := httptest.NewRecorder()

		gin.ServeHTTP(ideaRec, ideaReq)
		assertEqual(http.StatusCreated, ideaRec.Code, t, "Idea should have been created")

		idea := &domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(idea)

		assertUnequal(0, len(idea.Videos), t, "No videos found")
		assertUnequal(0, len(idea.Blogs), t, "No blogs found")
	})

	t.Run("ExistingVideo", func(t *testing.T) {
		db.Create(&domain.Video{
			Url:        "https://some-random-url.com",
			Identifier: "https://some-random-url.com",
		})

		ideaReqBody := []byte(`{
			"content": "Some random idea that I'd like to publish",
			"videos": [{"url": "https://some-random-url.com", "identifier": "https://some-random-url.com"}]
		}`)

		anotherIdeaReqBody := []byte(`{
			"content": "Yet another idea I'd like to publish",
			"videos": [{"url": "https://some-random-url.com", "identifier": "https://some-random-url.com"}]
		}`)

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

		gin.ServeHTTP(ideaRec, ideaReq)
		assertEqual(http.StatusCreated, ideaRec.Code, t, "Idea should have been created")

		gin.ServeHTTP(anotherIdeaRec, anotherIdeaReq)
		assertEqual(http.StatusCreated, anotherIdeaRec.Code, t, "Idea should have been created")

		currentVideosCount := len(fetchResources(db, []domain.Video{}))
		assertEqual(previousVideosCount, currentVideosCount, t, "Number of blogs increased")
	})

	t.Run("AnchorIdeas", func(t *testing.T) {
		ideas := seedIdeas(db, user)
		req := `{"content": "New idea with two anchor ideas", "anchors": [{"id": %d}, {"id": %d}]}`
		ideaReqBody := []byte(fmt.Sprintf(req, ideas[0].ID, ideas[1].ID))

		ideaReq, err := http.NewRequest(http.MethodPost, "/v1/ideas", bytes.NewReader(ideaReqBody))
		if err != nil {
			t.Fatalf("could not create request: %v\n", err)
		}

		ideaReq.Header.Add("Content-Type", "application/json")
		ideaReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		ideaRec := httptest.NewRecorder()

		previousIdeas := []domain.Idea{}
		db.Find(&previousIdeas)

		gin.ServeHTTP(ideaRec, ideaReq)
		assertEqual(http.StatusCreated, ideaRec.Code, t, "Idea should have been created")

		currentIdeas := []domain.Idea{}
		db.Find(&currentIdeas)

		ideaResp := &domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(ideaResp)
		assertEqual(len(ideas), len(ideaResp.Anchors), t, "Different number of anchor ideas")
		assertEqual(len(previousIdeas)+1, len(currentIdeas), t, "Different number of total ideas")
	})

	t.Run("BookWithChapter", func(t *testing.T) {
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
				"chapter": "2 - the greatest chapter"
			},
			{
				"url": "https://openlibrary.org/works/OL20984100F",
				"open_library_key": "OL20984100F",
				"title": "The Whatever Book Title",
				"year": 2000,
				"number_of_pages": 200,
				"open_library_id": 10320100,
				"language": "eng",
				"authors": [{"open_library_key": "OL1005931M", "full_name": "Best Writer"}],
				"chapter": "3 - another good chapter"
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
		assertEqual(http.StatusCreated, ideaRec.Code, t, "Idea should have been created")

		ideaResp := &domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(ideaResp)

		firstBookIdeaRel := &domain.BooksIdeas{
			IdeaID: ideaResp.ID,
			BookID: ideaResp.Books[0].ID,
		}
		db.Find(firstBookIdeaRel)

		secondBookIdeaRel := &domain.BooksIdeas{
			IdeaID: ideaResp.ID,
			BookID: ideaResp.Books[1].ID,
		}
		db.Find(secondBookIdeaRel)

		assertEqual("2 - the greatest chapter", firstBookIdeaRel.Chapter, t,
			fmt.Sprintf("Wrong chapter - book ID %d", ideaResp.Books[0].ID))

		assertEqual("3 - another good chapter", secondBookIdeaRel.Chapter, t,
			fmt.Sprintf("Wrong chapter - book ID %d", ideaResp.Books[0].ID))
	})

	t.Run("SameBook", func(t *testing.T) {
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
		assertEqual(http.StatusCreated, firstIdeaRec.Code, t, "First idea should have been created")

		ideaResp := &domain.Idea{}
		json.NewDecoder(firstIdeaRec.Body).Decode(ideaResp)
		bookID := ideaResp.Books[0].ID

		secondIdea := fmt.Sprintf((`{"content": "Another idea - book one", "books": [%s]}`), bookResource)
		secondIdeaReq, _ := http.NewRequest(http.MethodPost, "/v1/ideas", bytes.NewReader([]byte(secondIdea)))
		secondIdeaReq.Header.Add("Content-Type", "application/json")
		secondIdeaReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		secondIdeaRec := httptest.NewRecorder()

		gin.ServeHTTP(secondIdeaRec, secondIdeaReq)
		assertEqual(http.StatusCreated, secondIdeaRec.Code, t, "Second idea should have been created")

		bookIdeasRels := []domain.BooksIdeas{}
		db.Where([]domain.BooksIdeas{{BookID: bookID}}).Find(&bookIdeasRels)
		assertEqual(2, len(bookIdeasRels), t, fmt.Sprintf("Wrong amount of entries for book ID: %d", bookID))
	})

	t.Run("SameMovie", func(t *testing.T) {
		movieResource := `{
			"identifier": 62,
			"title": "2001: A Space Odyssey",
			"original_title": "2001: A Space Odyssey",
			"poster_path": "/15FumSExI9SRoL7QJWZAsA0b10c.jpg",
			"release_date": "1968-04-02",
			"runtime": 149,
			"original_language": "eng",
			"genres": [
				{"name": "Science Fiction", "name": "Mystery", "name": "Adventure"}
			]
		}`

		firstIdea := fmt.Sprintf((`{"content": "Some idea - movie one", "movies": [%s]}`), movieResource)
		firstIdeaReq, _ := http.NewRequest(http.MethodPost, "/v1/ideas", bytes.NewReader([]byte(firstIdea)))
		firstIdeaReq.Header.Add("Content-Type", "application/json")
		firstIdeaReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		firstIdeaRec := httptest.NewRecorder()

		gin.ServeHTTP(firstIdeaRec, firstIdeaReq)
		assertEqual(http.StatusCreated, firstIdeaRec.Code, t, "First idea should have been created")

		ideaResp := &domain.Idea{}
		json.NewDecoder(firstIdeaRec.Body).Decode(ideaResp)
		movieID := ideaResp.Movies[0].ID

		secondIdea := fmt.Sprintf((`{"content": "Another idea - movie one", "movies": [%s]}`), movieResource)
		secondIdeaReq, _ := http.NewRequest(http.MethodPost, "/v1/ideas", bytes.NewReader([]byte(secondIdea)))
		secondIdeaReq.Header.Add("Content-Type", "application/json")
		secondIdeaReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		secondIdeaRec := httptest.NewRecorder()

		gin.ServeHTTP(secondIdeaRec, secondIdeaReq)
		assertEqual(http.StatusCreated, secondIdeaRec.Code, t, "Second idea should have been created")

		moviesIdeasRels := []domain.IdeasMovies{}
		db.Where([]domain.IdeasMovies{{MovieID: movieID}}).Find(&moviesIdeasRels)
		assertEqual(2, len(moviesIdeasRels), t, fmt.Sprintf("Wrong amount of entries for movie ID: %d", movieID))
	})

	t.Run("MovieWithScene", func(t *testing.T) {
		scene := "When Joe meets Annie for the first time outside of the blue building."

		movieResource := fmt.Sprintf(`{
			"identifier": 10123,
			"title": "2002: A Galaxy Odyssey",
			"original_title": "2002: A Galaxy Odyssey",
			"poster_path": "/19XumSExI9SRoL7QJWZAsA0b10c.jpg",
			"release_date": "1970-01-01",
			"runtime": 180,
			"original_language": "eng",
			"genres": [
				{"name": "Science Fiction", "name": "Mystery", "name": "Adventure"}
			],
			"scene": "%s"
		}`, scene)

		idea := fmt.Sprintf((`{"content": "Some random idea", "movies": [%s]}`), movieResource)
		ideaReq, _ := http.NewRequest(http.MethodPost, "/v1/ideas", bytes.NewReader([]byte(idea)))
		ideaReq.Header.Add("Content-Type", "application/json")
		ideaReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authTokens.AccessToken))
		ideaRec := httptest.NewRecorder()

		gin.ServeHTTP(ideaRec, ideaReq)
		assertEqual(http.StatusCreated, ideaRec.Code, t, "Idea should have been created")

		ideaResp := domain.Idea{}
		json.NewDecoder(ideaRec.Body).Decode(&ideaResp)

		ideaMovieRel := domain.IdeasMovies{}
		db.Where(&domain.IdeasMovies{MovieID: ideaResp.Movies[0].ID, IdeaID: ideaResp.ID}).First(&ideaMovieRel)

		assertEqual(scene, ideaMovieRel.Scene, t, "Scene field not assigned correctly on ideas_movies")
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
	assertEqual(ideaRec.Code, http.StatusAccepted, t, "Idea should have been deleted")

	var idea domain.Idea
	db.First(&idea, ideas[1].ID)

	assertUnequal(0, idea.ID, t, "Idea should have been deleted")

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

	movie := domain.Movie{
		Identifier:       311,
		Title:            "Once Upon a Time in America",
		OriginalTitle:    "Once Upon a Time in America",
		PosterPath:       "/uPYa165sraN2c8gZBM9C47g3JoU.jpg",
		ReleaseDate:      "1984-05-23",
		Runtime:          229,
		OriginalLanguage: "en",
		Genres:           []domain.CinematicGenre{{Name: "Drama"}, {Name: "Crime"}},
	}

	fullIdea := &domain.Idea{
		UserID:  user.ID,
		Content: "Content for an idea anchored upon a blog",
		Blogs:   []domain.Blog{{Url: "https://some-blog.com", Category: "science"}},
		Videos:  []domain.Video{{Url: "https://some-youtube-video.com", YoutubeChannel: "cool-channel"}},
		Anchors: []*domain.Idea{emptyIdea},
		Books:   []domain.Book{book},
		Movies:  []domain.Movie{movie},
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
	db.Exec("delete from ideas_movies")
	db.Exec("delete from cinematic_genres_movies")
	db.Exec("delete from cinematic_genres")
	db.Exec("delete from movies")
	db.Exec("delete from authors")
	db.Exec("delete from books")
	db.Exec("delete from videos")
	db.Exec("delete from blogs")
	db.Exec("delete from ideas")
}

func checkIdeaAssociations(t *testing.T, idea *domain.Idea) {
	assertUnequal(0, len(idea.Blogs), t, "Blogs missing")
	assertUnequal(0, len(idea.Videos), t, "Videos missing")
	assertUnequal(0, len(idea.Anchors), t, "Anchors missing")
	assertUnequal(0, len(idea.Movies), t, "Movies missing")
	assertUnequal(0, len(idea.Movies[0].Genres), t, "Cinematic genres missing")
	assertUnequal(0, len(idea.Books), t, "Books missing")
	assertUnequal(0, len(idea.Books[0].Authors), t, "Authors missing")
}
