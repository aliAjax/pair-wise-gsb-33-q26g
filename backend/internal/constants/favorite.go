package constants

// FavoriteTargetType enumerates favorite target types.
const (
	FavoriteTargetPlant   = "plant"
	FavoriteTargetArticle = "article"
)

// ValidFavoriteTargets returns all accepted target types.
func ValidFavoriteTargets() []string {
	return []string{FavoriteTargetPlant, FavoriteTargetArticle}
}

// IsValidFavoriteTarget reports whether the target type is known.
func IsValidFavoriteTarget(t string) bool {
	for _, v := range ValidFavoriteTargets() {
		if v == t {
			return true
		}
	}
	return false
}
