package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/robiuzzaman4/donor-registry-backend/internal/domain"
	"github.com/robiuzzaman4/donor-registry-backend/internal/user"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) user.UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (
			name, phone, password, blood_group, role, gender, 
			date_of_birth, zila, upazila, local_address, is_available
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		u.Name, u.Phone, u.Password, u.BloodGroup, u.Role, u.Gender,
		u.DateOfBirth, u.Zila, u.Upazila, u.LocalAddress, u.IsAvailable,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("Failed to create user: %w", err)
	}

	return u, nil
}

func (r *userRepo) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	var u domain.User
	var lastDonated pgtype.Timestamptz
	query := `
		SELECT id, name, phone, password, blood_group, role, gender, date_of_birth,
		       zila, upazila, local_address, total_donate_count, is_verified, 
		       is_available, last_donated_at, created_at, updated_at
		FROM users 
		WHERE id = $1 AND is_deleted = false LIMIT 1`

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&u.ID, &u.Name, &u.Phone, &u.Password, &u.BloodGroup, &u.Role, &u.Gender, &u.DateOfBirth,
		&u.Zila, &u.Upazila, &u.LocalAddress,
		&u.TotalDonateCount, &u.IsVerified, &u.IsAvailable, &lastDonated, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to get user by id: %w", err)
	}
	if lastDonated.Valid {
		u.LastDonatedAt = lastDonated.Time
	}
	return &u, nil
}

func (r *userRepo) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var u domain.User
	var lastDonated pgtype.Timestamptz
	query := `
		SELECT id, name, phone, password, blood_group, role, gender, 
		       zila, upazila, local_address, total_donate_count, is_available, 
		       last_donated_at, created_at, updated_at
		FROM users WHERE phone = $1 AND is_deleted = false LIMIT 1`

	err := r.db.QueryRow(ctx, query, phone).Scan(
		&u.ID, &u.Name, &u.Phone, &u.Password, &u.BloodGroup, &u.Role, &u.Gender,
		&u.Zila, &u.Upazila, &u.LocalAddress,
		&u.TotalDonateCount, &u.IsAvailable, &lastDonated, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to get user by phone: %w", err)
	}
	if lastDonated.Valid {
		u.LastDonatedAt = lastDonated.Time
	}
	return &u, nil
}

func (r *userRepo) List(ctx context.Context, page, limit int) ([]*domain.User, int64, error) {
	// get total count for pagination metadata
	var total int64
	countQuery := `SELECT COUNT(id) FROM users WHERE is_deleted = false`
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to count users: %w", err)
	}

	// calculate offset internally
	offset := (page - 1) * limit

	// execute paginated query
	query := `
		SELECT id, name, phone, blood_group, gender, zila, upazila, local_address, 
		       total_donate_count, is_available, last_donated_at
		FROM users 
		WHERE is_deleted = false
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		var lastDonated pgtype.Timestamptz
		err := rows.Scan(
			&u.ID, &u.Name, &u.Phone, &u.BloodGroup, &u.Gender,
			&u.Zila, &u.Upazila, &u.LocalAddress,
			&u.TotalDonateCount, &u.IsAvailable, &lastDonated,
		)
		if err != nil {
			return nil, 0, err
		}
		if lastDonated.Valid {
			u.LastDonatedAt = lastDonated.Time
		}
		users = append(users, &u)
	}

	return users, total, nil
}

func (r *userRepo) FindByPhoneAndPassword(ctx context.Context, phone string, password string) (*domain.User, error) {
	var u domain.User
	query := `
		SELECT id, name, phone, password, role 
		FROM users 
		WHERE phone = $1 AND password = $2 AND is_deleted = false LIMIT 1`

	err := r.db.QueryRow(ctx, query, phone, password).Scan(
		&u.ID, &u.Name, &u.Phone, &u.Password, &u.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Login lookup failed: %w", err)
	}
	return &u, nil
}

func (r *userRepo) Update(ctx context.Context, userID string, u *domain.User) error {
	query := `
		UPDATE users 
		SET name = $1, 
		    blood_group = $2, 
		    gender = $3, 
		    zila = $4, 
		    upazila = $5, 
		    local_address = $6, 
		    is_available = $7, 
		    updated_at = NOW()
		WHERE id = $8 AND is_deleted = false`

	result, err := r.db.Exec(ctx, query,
		u.Name, u.BloodGroup, u.Gender, u.Zila, u.Upazila,
		u.LocalAddress, u.IsAvailable, userID,
	)

	if err != nil {
		return fmt.Errorf("Failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("No user found to update")
	}

	return nil
}
