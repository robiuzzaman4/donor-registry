package user

import (
	"context"

	"github.com/robiuzzaman4/donor-registry-backend/internal/domain"
)

type Service interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, userID string) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	List(ctx context.Context) ([]*domain.User, error)
	FindByPhoneAndPassword(ctx context.Context, phone string, password string) (*domain.User, error)
	Update(ctx context.Context, userID string, user *domain.User) error
}
