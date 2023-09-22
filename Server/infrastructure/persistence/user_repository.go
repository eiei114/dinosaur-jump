package infrastructure

import (
	"context"
	"example.com/domain"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	Conn *bun.DB
}

func NewUserRepository(Conn *bun.DB) *UserRepository {
	return &UserRepository{Conn: Conn}
}

func (u *UserRepository) AddUser(ctx context.Context, id, authToken, name string) error {
	user := &domain.User{
		Id:        id,
		AuthToken: authToken,
		Name:      name,
		HighScore: 0,
	}
	_, err := u.Conn.NewInsert().Model(user).Exec(ctx)
	return err
}

func (u *UserRepository) DeleteUser(ctx context.Context, id string) error {
	return nil
}

func (u *UserRepository) GetUserByUserId(ctx context.Context, id string) (*domain.User, error) {
	return nil, nil
}

func (u *UserRepository) GetUserByAuthToken(ctx context.Context, authToken string) (*domain.User, error) {
	user := new(domain.User)
	err := u.Conn.NewSelect().Model(user).Where("auth_token = ?", authToken).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) GetUserRanking(ctx context.Context) ([]*domain.UserRanking, error) {
	var users []domain.User

	err := u.Conn.NewSelect().
		Column("name", "high_score").
		Model(&users).
		OrderExpr("high_score DESC"). // ハイスコアで降順にソート
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	var userRankings []*domain.UserRanking
	for _, user := range users {
		ranking := &domain.UserRanking{
			Name:      user.Name,
			HighScore: user.HighScore,
		}
		userRankings = append(userRankings, ranking)
	}

	return userRankings, nil
}
