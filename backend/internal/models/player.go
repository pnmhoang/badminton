package models

// Player represents legacy player model (will be deprecated in favor of User)
type Player struct {
	BaseModel
	Name    string `json:"name" gorm:"not null"`
	Email   string `json:"email" gorm:"unique;not null"`
	Ranking int    `json:"ranking" gorm:"default:0"`
}
