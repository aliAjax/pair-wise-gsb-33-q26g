package model

import "time"

// DiseasePest is a plant disease or pest entry in the diagnostic manual.
type DiseasePest struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	PlantSpeciesID     uint      `gorm:"index" json:"plant_species_id"`
	Name               string    `gorm:"size:128;not null" json:"name"`
	Symptoms           string    `gorm:"type:text" json:"symptoms"`
	Cause              string    `gorm:"type:text" json:"cause"`
	Treatment          string    `gorm:"type:text" json:"treatment"`
	RecommendedMedicine string   `gorm:"size:255" json:"recommended_medicine"`
	Images             string    `gorm:"type:json" json:"images"`
	Keywords           string    `gorm:"size:255;index" json:"keywords"`
	CreatedAt          time.Time `json:"created_at"`
}
