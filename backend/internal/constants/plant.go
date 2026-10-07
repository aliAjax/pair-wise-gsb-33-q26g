package constants

// PlantType enumerates plant categories used across the whole stack.
const (
	PlantTypeFlower    = "flower"    // 观花
	PlantTypeFoliage   = "foliage"   // 观叶
	PlantTypeSucculent = "succulent" // 多肉
	PlantTypeAquatic   = "aquatic"   // 水生
)

// ValidPlantTypes returns all accepted plant type values.
func ValidPlantTypes() []string {
	return []string{PlantTypeFlower, PlantTypeFoliage, PlantTypeSucculent, PlantTypeAquatic}
}

// IsValidPlantType reports whether the given type is a known plant type.
func IsValidPlantType(t string) bool {
	for _, v := range ValidPlantTypes() {
		if v == t {
			return true
		}
	}
	return false
}
