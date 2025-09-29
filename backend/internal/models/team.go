package models

// Team represents a team for doubles tournaments
type Team struct {
	BaseModel
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`

	// Relations
	Players     []TeamPlayer     `json:"players,omitempty" gorm:"foreignKey:TeamID"`
	Tournaments []TournamentTeam `json:"tournaments,omitempty" gorm:"foreignKey:TeamID"`
}

// TeamPlayer represents the relationship between teams and players
type TeamPlayer struct {
	BaseModel
	TeamID   uint   `json:"team_id" gorm:"not null"`
	PlayerID uint   `json:"player_id" gorm:"not null"`
	Role     string `json:"role" gorm:"default:'player'"` // captain, player

	// Relations
	Team   Team `json:"team" gorm:"foreignKey:TeamID"`
	Player User `json:"player" gorm:"foreignKey:PlayerID"`
}

// TournamentPlayer represents player registration in tournaments
type TournamentPlayer struct {
	BaseModel
	TournamentID uint   `json:"tournament_id" gorm:"not null"`
	PlayerID     uint   `json:"player_id" gorm:"not null"`
	Status       string `json:"status" gorm:"default:'registered'"` // registered, confirmed, withdrawn

	// Relations
	Tournament Tournament `json:"tournament" gorm:"foreignKey:TournamentID"`
	Player     User       `json:"player" gorm:"foreignKey:PlayerID"`
}

// TournamentTeam represents team registration in tournaments
type TournamentTeam struct {
	BaseModel
	TournamentID uint   `json:"tournament_id" gorm:"not null"`
	TeamID       uint   `json:"team_id" gorm:"not null"`
	Status       string `json:"status" gorm:"default:'registered'"` // registered, confirmed, withdrawn

	// Relations
	Tournament Tournament `json:"tournament" gorm:"foreignKey:TournamentID"`
	Team       Team       `json:"team" gorm:"foreignKey:TeamID"`
}
