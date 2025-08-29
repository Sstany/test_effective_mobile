package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"go.uber.org/zap"

	repo "subscription-service/internal/adapter/repo/mock"
	"subscription-service/internal/app/entity"
	"subscription-service/internal/app/usecase"
	pkg "subscription-service/internal/pkg/utils"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subscriptionRepo := repo.NewMockSubscriptionRepo(ctrl)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	subscriptionUsecase, err := usecase.NewSubscription(subscriptionRepo, nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	createRequest := entity.CreateSubscriptionRequest{
		Title:     "Premium",
		Price:     1000,
		UserID:    "user123",
		StartDate: time.Now(),
		EndDate:   pkg.PointerTo(time.Now().Add(time.Hour * 24 * 30)),
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}

	ctx := context.Background()

	subscriptionRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context,
			req entity.CreateSubscriptionRequest,
		) error {
			if req.ID == "" {
				return errors.New("ID should be generated")
			}
			return nil
		})

	result, err := subscriptionUsecase.Create(ctx, createRequest)
	if err != nil {
		t.Error(err)
	}

	if result.ID == "" {
		t.Error("expected ID to be generated")
	}
	if result.Title != createRequest.Title {
		t.Errorf("expected title %s, got %s", createRequest.Title, result.Title)
	}
}

func TestCreateWithRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subscriptionRepo := repo.NewMockSubscriptionRepo(ctrl)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	subscriptionUsecase, err := usecase.NewSubscription(subscriptionRepo, nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	createRequest := entity.CreateSubscriptionRequest{
		Title:     "Premium",
		Price:     1000,
		UserID:    "user123",
		StartDate: time.Now(),
		EndDate:   pkg.PointerTo(time.Now().Add(time.Hour * 24 * 30)),
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}

	ctx := context.Background()

	expectedErr := errors.New("database error")
	subscriptionRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedErr)

	_, err = subscriptionUsecase.Create(ctx, createRequest)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subscriptionRepo := repo.NewMockSubscriptionRepo(ctrl)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	subscriptionUsecase, err := usecase.NewSubscription(subscriptionRepo, nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	subscriptionID := uuid.NewString()
	expectedSubscription := &entity.Subscription{
		ID:     subscriptionID,
		Title:  "Premium",
		Price:  1000,
		UserID: "user123",
	}

	ctx := context.Background()

	subscriptionRepo.EXPECT().GetSubscription(ctx, subscriptionID).Return(expectedSubscription, nil)

	result, err := subscriptionUsecase.Read(ctx, subscriptionID)
	if err != nil {
		t.Error(err)
	}

	if result.ID != expectedSubscription.ID {
		t.Errorf("expected ID %s, got %s", expectedSubscription.ID, result.ID)
	}
}

func TestReadWithNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subscriptionRepo := repo.NewMockSubscriptionRepo(ctrl)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	subscriptionUsecase, err := usecase.NewSubscription(subscriptionRepo, nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	subscriptionID := uuid.NewString()

	ctx := context.Background()

	expectedErr := errors.New("not found")
	subscriptionRepo.EXPECT().GetSubscription(ctx, subscriptionID).Return(nil, expectedErr)

	_, err = subscriptionUsecase.Read(ctx, subscriptionID)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subscriptionRepo := repo.NewMockSubscriptionRepo(ctrl)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	subscriptionUsecase, err := usecase.NewSubscription(subscriptionRepo, nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	updateRequest := entity.UpdateSubscriptionRequest{
		ID:    uuid.NewString(),
		Title: "Updated Premium",
		Price: 1500,
	}

	ctx := context.Background()

	subscriptionRepo.EXPECT().Update(ctx, updateRequest).Return(nil)

	err = subscriptionUsecase.Update(ctx, updateRequest)
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subscriptionRepo := repo.NewMockSubscriptionRepo(ctrl)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	subscriptionUsecase, err := usecase.NewSubscription(subscriptionRepo, nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	subscriptionID := uuid.NewString()

	ctx := context.Background()

	subscriptionRepo.EXPECT().Delete(ctx, subscriptionID).Return(nil)

	err = subscriptionUsecase.Delete(ctx, subscriptionID)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteWithRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subscriptionRepo := repo.NewMockSubscriptionRepo(ctrl)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	subscriptionUsecase, err := usecase.NewSubscription(subscriptionRepo, nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	subscriptionID := uuid.NewString()

	ctx := context.Background()

	expectedErr := errors.New("delete failed")
	subscriptionRepo.EXPECT().Delete(ctx, subscriptionID).Return(expectedErr)

	err = subscriptionUsecase.Delete(ctx, subscriptionID)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subscriptionRepo := repo.NewMockSubscriptionRepo(ctrl)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	subscriptionUsecase, err := usecase.NewSubscription(subscriptionRepo, nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	filter := entity.ListSubscriptionFilter{
		UserID: pkg.PointerTo("60601fee-2bf1-4721-ae6f-7636e79a0cba"),
		Limit:  pkg.PointerTo(10),
		Offset: pkg.PointerTo(0),
	}

	expectedSubscriptions := []entity.Subscription{
		{
			ID:     uuid.NewString(),
			Title:  "Premium",
			Price:  1000,
			UserID: "user123",
		},
		{
			ID:     uuid.NewString(),
			Title:  "Basic",
			Price:  500,
			UserID: "user123",
		},
	}

	ctx := context.Background()

	subscriptionRepo.EXPECT().List(ctx, filter).Return(expectedSubscriptions, nil)

	result, err := subscriptionUsecase.List(ctx, filter)
	if err != nil {
		t.Error(err)
	}

	if len(result) != len(expectedSubscriptions) {
		t.Errorf("expected %d subscriptions, got %d", len(expectedSubscriptions), len(result))
	}
}

func TestListWithEmptyResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subscriptionRepo := repo.NewMockSubscriptionRepo(ctrl)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	subscriptionUsecase, err := usecase.NewSubscription(subscriptionRepo, nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	filter := entity.ListSubscriptionFilter{
		UserID: pkg.PointerTo("60601fee-2bf1-4721-ae6f-7636e79a0cba"),
		Limit:  pkg.PointerTo(10),
		Offset: pkg.PointerTo(0),
	}

	ctx := context.Background()

	subscriptionRepo.EXPECT().List(ctx, filter).Return([]entity.Subscription{}, nil)

	result, err := subscriptionUsecase.List(ctx, filter)
	if err != nil {
		t.Error(err)
	}

	if len(result) != 0 {
		t.Errorf("expected empty result, got %d items", len(result))
	}
}

func TestSum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subscriptionRepo := repo.NewMockSubscriptionRepo(ctrl)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	subscriptionUsecase, err := usecase.NewSubscription(subscriptionRepo, nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	filter := entity.ListSubscriptionFilter{
		UserID: pkg.PointerTo("60601fee-2bf1-4721-ae6f-7636e79a0cba"),
	}

	expectedSum := int64(2500)

	ctx := context.Background()

	subscriptionRepo.EXPECT().Sum(ctx, filter).Return(expectedSum, nil)

	result, err := subscriptionUsecase.Sum(ctx, filter)
	if err != nil {
		t.Error(err)
	}

	if result != expectedSum {
		t.Errorf("expected sum %d, got %d", expectedSum, result)
	}
}

func TestSumWithZeroResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subscriptionRepo := repo.NewMockSubscriptionRepo(ctrl)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	subscriptionUsecase, err := usecase.NewSubscription(subscriptionRepo, nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	filter := entity.ListSubscriptionFilter{
		UserID: pkg.PointerTo("60601fee-2bf1-4721-ae6f-7636e79a0cba"),
	}

	ctx := context.Background()

	subscriptionRepo.EXPECT().Sum(ctx, filter).Return(int64(0), nil)

	result, err := subscriptionUsecase.Sum(ctx, filter)
	if err != nil {
		t.Error(err)
	}

	if result != 0 {
		t.Errorf("expected sum 0, got %d", result)
	}
}

func TestSumWithRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subscriptionRepo := repo.NewMockSubscriptionRepo(ctrl)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	subscriptionUsecase, err := usecase.NewSubscription(subscriptionRepo, nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	filter := entity.ListSubscriptionFilter{
		UserID: pkg.PointerTo("60601fee-2bf1-4721-ae6f-7636e79a0cba"),
	}

	ctx := context.Background()

	expectedErr := errors.New("sum calculation failed")
	subscriptionRepo.EXPECT().Sum(ctx, filter).Return(int64(0), expectedErr)

	_, err = subscriptionUsecase.Sum(ctx, filter)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
