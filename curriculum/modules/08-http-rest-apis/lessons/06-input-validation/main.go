package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var validEmail = regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
var alphanumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
var validRoles = map[string]bool{
	"admin": true, "editor": true, "viewer": true,
}

func validateUser(input map[string]interface{}) []ValidationError {
	var errs []ValidationError

	username, _ := input["username"].(string)
	username = strings.TrimSpace(username)
	if username == "" {
		errs = append(errs, ValidationError{Field: "username", Message: "username is required"})
	} else {
		if len(username) < 3 || len(username) > 30 {
			errs = append(errs, ValidationError{Field: "username", Message: "username must be 3-30 characters"})
		}
		if !alphanumeric.MatchString(username) {
			errs = append(errs, ValidationError{Field: "username", Message: "username must be alphanumeric"})
		}
	}

	ageVal, ageOk := input["age"]
	if !ageOk {
		errs = append(errs, ValidationError{Field: "age", Message: "age is required"})
	} else {
		ageFloat, ok := ageVal.(float64)
		if !ok {
			errs = append(errs, ValidationError{Field: "age", Message: "age must be a number"})
		} else {
			age := int(ageFloat)
			if age < 13 || age > 150 {
				errs = append(errs, ValidationError{Field: "age", Message: "age must be between 13 and 150"})
			}
		}
	}

	if email, ok := input["email"].(string); ok && email != "" {
		if !validEmail.MatchString(email) {
			errs = append(errs, ValidationError{Field: "email", Message: "invalid email format"})
		}
	}

	if role, ok := input["role"].(string); ok && role != "" {
		if !validRoles[role] {
			errs = append(errs, ValidationError{Field: "role", Message: fmt.Sprintf("role must be one of: admin, editor, viewer")})
		}
	}

	return errs
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}

	errs := validateUser(input)
	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": errs})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
