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
	BloodGroup       UserBloodGroup `json:"bloodGroup" db:"blood_group"`
	Role             UserRole       `json:"role" db:"role"`
	Gender           UserGender     `json:"gender" db:"gender"`
	DateOfBirth      time.Time      `json:"dateOfBirth" db:"date_of_birth"`
	Zila             string         `json:"zila" db:"zila"`
	Upazila          string         `json:"upazila" db:"upazila"`
	LocalAddress     string         `json:"localAddress" db:"local_address"`
	TotalDonateCount int            `json:"totalDonateCount" db:"total_donate_count"`
	IsVerified       bool           `json:"isVerified" db:"is_verified"`
	IsAvailable      bool           `json:"isAvailable" db:"is_available"`
	IsDeleted        bool           `json:"IsDeleted" db:"is_deleted"`
	LastDonatedAt    time.Time      `json:"lastDonatedAt" db:"last_donated_at"`
	CreatedAt        time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time      `json:"updatedAt" db:"updated_at"`
}
