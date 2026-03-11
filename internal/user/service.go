package user

import (
	"context"

	"github.com/robiuzzaman4/donor-registry-backend/internal/domain"
)

type service struct {
	ctx      context.Context
	userRepo UserRepo
}

func NewService(ctx context.Context, userRepo UserRepo) Service {
	return &service{
		ctx:      ctx,
		userRepo: userRepo,
	}
}

func (svc service) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	return svc.userRepo.Create(ctx, user)
}

func (svc service) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	return svc.userRepo.GetByID(ctx, userID)
}

func (svc service) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	return svc.userRepo.GetByPhone(ctx, phone)
}

func (s *service) List(ctx context.Context, page, limit int) ([]*domain.User, int64, error) {
	return s.userRepo.List(ctx, page, limit)
}

func (svc service) FindByPhoneAndPassword(ctx context.Context, phone string, password string) (*domain.User, error) {
	return svc.userRepo.FindByPhoneAndPassword(ctx, phone, password)
}

func (svc service) Update(ctx context.Context, userID string, user *domain.User) error {
	return svc.userRepo.Update(ctx, userID, user)
}
