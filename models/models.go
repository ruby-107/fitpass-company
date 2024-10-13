package models

// User Model
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Profile Model
type Profile struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	ProfileName string `json:"profile_name"`
}

// UserProfileResponse Model
type UserProfileResponse struct {
	ID          int    `json:"id"`
	UserID      User   `json:"user_id"`
	ProfileName string `json:"profile_name"`
}
