package domain

import (
	"time"

	"github.com/google/uuid"
)

// Plan represents subscription tiers
type Plan string

const (
	PlanFree     Plan = "free"
	PlanBasic    Plan = "basic"
	PlanPro      Plan = "pro"
	PlanBusiness Plan = "business"
)

// User represents a SaaS customer
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Plan         Plan      `json:"plan" db:"plan"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// NewUser creates a new user with defaults
func NewUser(email, passwordHash string) *User {
	now := time.Now()
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		Plan:         PlanFree,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// PlanLimits returns the account limits for each plan
func (p Plan) Limits() PlanLimits {
	switch p {
	case PlanBasic:
		return PlanLimits{MaxAccounts: 3, ScanIntervalHours: 12}
	case PlanPro:
		return PlanLimits{MaxAccounts: 10, ScanIntervalHours: 6}
	case PlanBusiness:
		return PlanLimits{MaxAccounts: 50, ScanIntervalHours: 1}
	default: // Free
		return PlanLimits{MaxAccounts: 1, ScanIntervalHours: 24}
	}
}

type PlanLimits struct {
	MaxAccounts       int
	ScanIntervalHours int
}
