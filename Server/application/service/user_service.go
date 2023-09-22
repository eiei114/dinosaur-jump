package service

import (
	"context"
	"example.com/domain"
	"example.com/domain/repository"
	"github.com/google/uuid"
)

type UserService struct {
	UserRepository repository.UserRepository
}

func NewUserService(UserRepository repository.UserRepository) *UserService {
	return &UserService{UserRepository}
}

func (u *UserService) Add(ctx context.Context, name string) (string, error) {
	// UUIDでユーザIDと認証トークンを生成
	userID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	authToken, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	err = u.UserRepository.AddUser(ctx, userID.String(), authToken.String(), name)
	if err != nil {
		return "", err
	}

	return authToken.String(), nil
}

func (u *UserService) Delete(ctx context.Context, id string) (string, error) {
	_ = u.UserRepository.DeleteUser(ctx, id)
	return "", nil
}

func (u *UserService) GetUserByUserId(ctx context.Context, id string) (*domain.User, error) {
	var user *domain.User
	user, _ = u.UserRepository.GetUserByUserId(ctx, id)
	return user, nil
}

func (u *UserService) GetUserByAuthToken(ctx context.Context, authToken string) (*domain.User, error) {
	var user *domain.User
	user, _ = u.UserRepository.GetUserByAuthToken(ctx, authToken)
	return user, nil
}

func (u *UserService) GetUserRanking(ctx context.Context) ([]*domain.UserRanking, error) {
	userRankings, err := u.UserRepository.GetUserRanking(ctx)
	if err != nil {
		return nil, err
	}
	return userRankings, nil
}
