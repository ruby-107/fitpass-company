package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"fitpass.com/database"
	"fitpass.com/models"
)

// Create User
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !isValidEmail(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Begin a transaction
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	sqlStatement := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
	err = tx.QueryRow(sqlStatement, user.Name, user.Email).Scan(&user.ID)

	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Transaction commit failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// Create Profile
func CreateProfile(w http.ResponseWriter, r *http.Request) {
	var profile models.Profile

	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Begin a transaction
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Transaction rollback failed: %v", rollbackErr)
			}
		}
	}()

	// Check if the user exists
	var user models.User
	err = tx.QueryRow(`SELECT id, name, email FROM users WHERE id=$1`, profile.UserID).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	// Check if the user already has a profile
	var existingProfileID int
	err = tx.QueryRow(`SELECT id FROM profiles WHERE user_id=$1`, profile.UserID).Scan(&existingProfileID)

	if err == nil {
		http.Error(w, "User already has a profile", http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sqlStatement := `INSERT INTO profiles (user_id, profile_name) VALUES ($1, $2) RETURNING id`
	err = tx.QueryRow(sqlStatement, profile.UserID, profile.ProfileName).Scan(&profile.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Commit the transaction after both operations are successful
	if err := tx.Commit(); err != nil {
		http.Error(w, "Transaction commit failed", http.StatusInternalServerError)
		return
	}

	response := models.UserProfileResponse{
		ID: profile.ID,
		UserID: models.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		ProfileName: profile.ProfileName,
	}

	json.NewEncoder(w).Encode(response)
}

// Email validation function
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
