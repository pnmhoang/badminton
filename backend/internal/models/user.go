package models

import "golang.org/x/crypto/bcrypt"

// UserRole defines user roles
type UserRole string

const (
	RolePlayer UserRole = "player"
	RoleAdmin  UserRole = "admin"
)

// User represents a system user (both players and admins)
type User struct {
	BaseModel
	Username string   `json:"username" gorm:"unique;not null"`
	Email    string   `json:"email" gorm:"unique;not null"`
	Password string   `json:"-" gorm:"not null"` // Hidden in JSON
	FullName string   `json:"full_name" gorm:"not null"`
	Role     UserRole `json:"role" gorm:"default:'player'"`
	IsActive bool     `json:"is_active" gorm:"default:true"`

	// Player-specific fields (only used when Role = player)
	Ranking int `json:"ranking,omitempty" gorm:"default:0"`

	// Relations
	PlayerTeams      []TeamPlayer       `json:"player_teams,omitempty" gorm:"foreignKey:PlayerID"`
	AdminTournaments []Tournament       `json:"admin_tournaments,omitempty" gorm:"foreignKey:AdminID"`
	Registrations    []TournamentPlayer `json:"registrations,omitempty" gorm:"foreignKey:PlayerID"`
}

// HashPassword hashes the user password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies the user password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsAdmin checks if user is admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsPlayer checks if user is player
func (u *User) IsPlayer() bool {
	return u.Role == RolePlayer
}
