package domain

import "time"

type UserBloodGroup string
type UserRole string
type UserGender string

const (
	APositive  UserBloodGroup = "A+"
	ANegative  UserBloodGroup = "A-"
	BPositive  UserBloodGroup = "B+"
	BNegative  UserBloodGroup = "B-"
	OPositive  UserBloodGroup = "O+"
	ONegative  UserBloodGroup = "O-"
	ABPositive UserBloodGroup = "AB+"
	ABNegative UserBloodGroup = "AB-"
)
const (
	RoleAdmin UserRole = "ADMIN"
	RoleUser  UserRole = "USER"
)
const (
	Male   UserGender = "MALE"
	Female UserGender = "FEMALE"
)

type User struct {
	ID               string         `json:"id" db:"id"`
	Name             string         `json:"name" db:"name"`
	Phone            string         `json:"phone" db:"phone"`
	Password         string         `json:"-" db:"password"`
	BloodGroup       UserBloodGroup `json:"blood_group" db:"blood_group"`
	Role             UserRole       `json:"role" db:"role"`
	Gender           UserGender     `json:"gender" db:"gender"`
	DateOfBirth      time.Time      `json:"date_of_birth" db:"date_of_birth"`
	Zila             string         `json:"zila" db:"zila"`
	Upazila          string         `json:"upazila" db:"upazila"`
	LocalAddress     string         `json:"local_address" db:"local_address"`
	TotalDonateCount int            `json:"total_donate_count" db:"total_donate_count"`
	IsVerified       bool           `json:"is_verified" db:"is_verified"`
	IsAvailable      bool           `json:"is_available" db:"is_available"`
	IsDeleted        bool           `json:"is_deleted" db:"is_deleted"`
	LastDonatedAt    time.Time      `json:"last_donated_at" db:"last_donated_at"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}
