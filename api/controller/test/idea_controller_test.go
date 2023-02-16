package controller

import "testing"

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
	// logged in user is needed
	t.Skip()
}

func TestDeleteIdeaByID(t *testing.T) {
	// logged in user is needed
	t.Skip()
}
