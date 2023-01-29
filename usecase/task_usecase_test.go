package usecase_test

// func TestFetchByUserID(t *testing.T) {
// 	mockTaskRepository := new(mocks.TaskRepository)
// 	userObjectID := primitive.NewObjectID()
// 	userID := userObjectID.Hex()

// 	t.Run("success", func(t *testing.T) {

// 		mockTask := domain.Task{
// 			Title:  "Test Title",
// 			UserID: uint,
// 		}

// 		mockListTask := make([]domain.Task, 0)
// 		mockListTask = append(mockListTask, mockTask)

// 		mockTaskRepository.On("FetchByUserID", mock.Anything, userID).Return(mockListTask, nil).Once()

// 		u := usecase.NewTaskUsecase(mockTaskRepository, time.Second*2)

// 		list, err := u.FetchByUserID(context.Background(), userID)

// 		assert.NoError(t, err)
// 		assert.NotNil(t, list)
// 		assert.Len(t, list, len(mockListTask))

// 		mockTaskRepository.AssertExpectations(t)
// 	})

// 	t.Run("error", func(t *testing.T) {
// 		mockTaskRepository.On("FetchByUserID", mock.Anything, userID).Return(nil, errors.New("Unexpected")).Once()

// 		u := usecase.NewTaskUsecase(mockTaskRepository, time.Second*2)

// 		list, err := u.FetchByUserID(context.Background(), userID)

// 		assert.Error(t, err)
// 		assert.Nil(t, list)

// 		mockTaskRepository.AssertExpectations(t)
// 	})

// }
