package dto

import "time"

// GardenAddRequest adds a plant to the user's garden.
type GardenAddRequest struct {
	PlantSpeciesID uint      `json:"plant_species_id" binding:"required"`
	Nickname       string    `json:"nickname" binding:"omitempty,max=64"`
	OwnedSince     time.Time `json:"owned_since"`
	Location       string    `json:"location" binding:"omitempty,max=128"`
}

// GardenBindRequest binds a reminder to a garden item.
type GardenBindRequest struct {
	ReminderID uint `json:"care_reminder_id" binding:"required"`
}
