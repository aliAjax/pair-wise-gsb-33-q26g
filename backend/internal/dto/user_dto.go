package dto

// RegisterRequest is the registration payload.
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
	Nickname string `json:"nickname" binding:"omitempty,max=64"`
}

// LoginRequest is the login payload.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest is the profile update payload.
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,max=64"`
	Bio      string `json:"bio" binding:"omitempty,max=512"`
	Avatar   string `json:"avatar" binding:"omitempty,max=255"`
}

// LoginResponse carries the token and user info.
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}
