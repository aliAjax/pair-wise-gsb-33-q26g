package dto

// PlantCreateRequest is the admin payload for creating a plant species.
type PlantCreateRequest struct {
	Family           string  `json:"family" binding:"required,max=64"`
	Genus            string  `json:"genus" binding:"required,max=64"`
	Name             string  `json:"name" binding:"required,max=128"`
	Alias            string  `json:"alias" binding:"omitempty,max=128"`
	Type             string  `json:"type" binding:"required"`
	Origin           string  `json:"origin" binding:"omitempty,max=128"`
	TempMin          float64 `json:"temp_min"`
	TempMax          float64 `json:"temp_max"`
	LightRequirement string  `json:"light_requirement" binding:"omitempty,max=255"`
	WaterFrequency   string  `json:"water_frequency" binding:"omitempty,max=255"`
	Description      string  `json:"description"`
	ImageURLs        string  `json:"image_urls"`
}
