package dto

// PestCreateRequest is the admin payload for disease/pest entries.
type PestCreateRequest struct {
	PlantSpeciesID      uint   `json:"plant_species_id"`
	Name                string `json:"name" binding:"required,max=128"`
	Symptoms            string `json:"symptoms"`
	Cause               string `json:"cause"`
	Treatment           string `json:"treatment"`
	RecommendedMedicine string `json:"recommended_medicine" binding:"omitempty,max=255"`
	Images              string `json:"images"`
	Keywords            string `json:"keywords" binding:"omitempty,max=255"`
}
