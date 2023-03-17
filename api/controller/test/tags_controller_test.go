package controller

import "testing"

func TestTags(t *testing.T) {
	gin, db := setup()
	_, user := signup(gin, db, sampleUser())

	t.Run("Create", func(t *testing.T) {
		t.Skip()
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
