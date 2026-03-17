package user

import (
	"context"

	"github.com/robiuzzaman4/donor-registry/internal/domain"
	userHandler "github.com/robiuzzaman4/donor-registry/internal/rest/handler/user"
)

type Service interface {
	userHandler.Service
}

type UserRepo interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, userID string) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	List(ctx context.Context, page, limit int) ([]*domain.User, int64, error)
	Update(ctx context.Context, userID string, user *domain.User) error
}
