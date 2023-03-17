package controller

import "testing"

func TestTags(t *testing.T) {
	gin, db := setup()
	_, user := signup(gin, db, sampleUser())

	t.Run("Create", func(t *testing.T) {
		t.Skip()
	})

	t.Cleanup(func() {
		cleanupDatabase(db)
		cleanupUser(db, user.Name)
	})
}
