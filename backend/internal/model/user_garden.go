package model

import "time"

// UserGarden represents a plant owned by a user inside their garden list.
type UserGarden struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"index:idx_garden_user_plant,unique;not null" json:"user_id"`
	PlantSpeciesID uint      `gorm:"index:idx_garden_user_plant,unique;not null" json:"plant_species_id"`
	Nickname       string    `gorm:"size:64" json:"nickname"`
	OwnedSince     time.Time `gorm:"type:date" json:"owned_since"`
	Location       string    `gorm:"size:128" json:"location"`
	CareReminderID uint      `json:"care_reminder_id"`
	CreatedAt      time.Time `json:"created_at"`
}
