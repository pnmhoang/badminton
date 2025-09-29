package models

import "time"

// TournamentType defines tournament types
type TournamentType string

const (
	TournamentSingles TournamentType = "singles"
	TournamentDoubles TournamentType = "doubles"
)

// TournamentStatus defines tournament status
type TournamentStatus string

const (
	TournamentUpcoming  TournamentStatus = "upcoming"
	TournamentOngoing   TournamentStatus = "ongoing"
	TournamentCompleted TournamentStatus = "completed"
	TournamentCancelled TournamentStatus = "cancelled"
)

// Tournament represents a badminton tournament
type Tournament struct {
	BaseModel
	Name        string           `json:"name" gorm:"not null"`
	Description string           `json:"description"`
	Type        TournamentType   `json:"type" gorm:"default:'singles'"` // singles or doubles
	StartDate   time.Time        `json:"start_date"`
	EndDate     time.Time        `json:"end_date"`
	Status      TournamentStatus `json:"status" gorm:"default:'upcoming'"`
	MaxPlayers  int              `json:"max_players" gorm:"default:16"`
	MaxTeams    int              `json:"max_teams" gorm:"default:8"` // for doubles tournaments
	EntryFee    float64          `json:"entry_fee" gorm:"default:0"`
	PrizePool   float64          `json:"prize_pool" gorm:"default:0"`
	AdminID     uint             `json:"admin_id" gorm:"not null"`

	// Relations
	Admin   User               `json:"admin" gorm:"foreignKey:AdminID"`
	Matches []Match            `json:"matches,omitempty" gorm:"foreignKey:TournamentID"`
	Players []TournamentPlayer `json:"players,omitempty" gorm:"foreignKey:TournamentID"`
	Teams   []TournamentTeam   `json:"teams,omitempty" gorm:"foreignKey:TournamentID"`
}

// IsTeamTournament checks if this is a team tournament
func (t *Tournament) IsTeamTournament() bool {
	return t.Type == TournamentDoubles
}

// GetMaxParticipants returns max participants based on tournament type
func (t *Tournament) GetMaxParticipants() int {
	if t.IsTeamTournament() {
		return t.MaxTeams
	}
	return t.MaxPlayers
}
