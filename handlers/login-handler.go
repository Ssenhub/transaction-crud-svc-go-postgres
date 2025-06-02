package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"transaction-crud-svc-go-postgres/models"

	"github.com/golang-jwt/jwt/v5"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.User

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	userName := []byte(os.Getenv("USER_NAME"))
	password := []byte(os.Getenv("PASSWORD"))

	if creds.UserName != string(userName) || creds.Password != string(password) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": creds.UserName,
		"exp":  time.Now().Add(time.Hour).Unix(),
	})

	jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))

	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Failed to sign token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}
