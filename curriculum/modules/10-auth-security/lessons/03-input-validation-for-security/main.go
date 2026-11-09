package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var (
	usernameRE = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
	docIDRE    = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)
)

type CreateUserRequest struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password,omitempty"`
	ReferralCode string `json:"referral_code,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Message)
}

func validateCreateUser(req CreateUserRequest) []ValidationError {
	var errs []ValidationError

	if !usernameRE.MatchString(req.Username) {
		errs = append(errs, ValidationError{
			Field: "username", Message: "must be 3-30 alphanumeric characters or underscores",
		})
	}

	if _, err := mailParseAddress(req.Email); err != nil {
		errs = append(errs, ValidationError{
			Field: "email", Message: "must be a valid email address",
		})
	}

	if len(req.Password) < 8 || len(req.Password) > 128 {
		errs = append(errs, ValidationError{
			Field: "password", Message: "must be between 8 and 128 characters",
		})
	}

	if req.ReferralCode != "" && !docIDRE.MatchString(req.ReferralCode) {
		errs = append(errs, ValidationError{
			Field: "referral_code", Message: "contains invalid characters",
		})
	}

	return errs
}

func mailParseAddress(addr string) (struct{}, error) {
	// Simple email validation: contains @ and a domain part.
	if !strings.Contains(addr, "@") {
		return struct{}{}, errors.New("invalid email")
	}
	parts := strings.Split(addr, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return struct{}{}, errors.New("invalid email")
	}
	if !strings.Contains(parts[1], ".") {
		return struct{}{}, errors.New("invalid email")
	}
	return struct{}{}, nil
}

func ValidateAPIKey(key string) error {
	if len(key) == 0 {
		return errors.New("key must not be empty")
	}
	if len(key) < 32 {
		return fmt.Errorf("key too short: %d < 32", len(key))
	}
	if len(key) > 128 {
		return fmt.Errorf("key too long: %d > 128", len(key))
	}
	for _, c := range key {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return fmt.Errorf("invalid character: %c", c)
		}
	}
	compromised := map[string]bool{
		"compromised-key-1234567890abcdef1234567890": true,
	}
	if compromised[key] {
		return errors.New("key has been compromised")
	}
	return nil
}

func validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name must not be empty")
	}
	if !usernameRE.MatchString(name) {
		return errors.New("name contains invalid characters")
	}
	return nil
}

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"invalid JSON: %s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		if errs := validateCreateUser(req); len(errs) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(map[string]interface{}{"errors": errs})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "created", "username": req.Username})
	})

	fmt.Println("Input validation server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
