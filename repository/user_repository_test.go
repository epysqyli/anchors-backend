package repository_test

// func TestCreate(t *testing.T) {

// 	var databaseHelper *mocks.Database
// 	var collectionHelper *mocks.Collection

// 	databaseHelper = &mocks.Database{}
// 	collectionHelper = &mocks.Collection{}

// 	collectionName := domain.CollectionUser

// 	mockUser := &domain.User{
// 		Name:     "Test",
// 		Email:    "test@gmail.com",
// 		Password: "password",
// 	}

// 	mockEmptyUser := &domain.User{}
// 	mockUserID := primitive.NewObjectID()

// 	t.Run("success", func(t *testing.T) {

// 		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("*domain.User")).Return(mockUserID, nil).Once()

// 		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

// 		ur := repository.NewUserRepository(databaseHelper)

// 		err := ur.Create(context.Background(), mockUser)

// 		assert.NoError(t, err)

// 		collectionHelper.AssertExpectations(t)
// 	})

// 	t.Run("error", func(t *testing.T) {
// 		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("*domain.User")).Return(mockEmptyUser, errors.New("Unexpected")).Once()

// 		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

// 		ur := repository.NewUserRepository(databaseHelper)

// 		err := ur.Create(context.Background(), mockEmptyUser)

// 		assert.Error(t, err)

// 		collectionHelper.AssertExpectations(t)
// 	})

// }
