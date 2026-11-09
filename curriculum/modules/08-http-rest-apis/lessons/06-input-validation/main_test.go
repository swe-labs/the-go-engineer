package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateUser_Valid(t *testing.T) {
	input := map[string]interface{}{
		"username": "alice",
		"age":      30.0,
		"email":    "alice@example.com",
		"role":     "admin",
	}
	errs := validateUser(input)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateUser_MissingUsername(t *testing.T) {
	input := map[string]interface{}{
		"age": 30.0,
	}
	errs := validateUser(input)
	if len(errs) == 0 {
		t.Fatal("expected errors, got none")
	}
	found := false
	for _, e := range errs {
		if e.Field == "username" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected username error, got %v", errs)
	}
}

func TestValidateUser_ShortUsername(t *testing.T) {
	input := map[string]interface{}{
		"username": "ab",
		"age":      25.0,
	}
	errs := validateUser(input)
	hasLenErr := false
	for _, e := range errs {
		if e.Field == "username" && strings.Contains(e.Message, "3-30") {
			hasLenErr = true
		}
	}
	if !hasLenErr {
		t.Errorf("expected username length error, got %v", errs)
	}
}

func TestValidateUser_InvalidAge(t *testing.T) {
	input := map[string]interface{}{
		"username": "bob",
		"age":      10.0,
	}
	errs := validateUser(input)
	hasAgeErr := false
	for _, e := range errs {
		if e.Field == "age" && strings.Contains(e.Message, "13 and 150") {
			hasAgeErr = true
		}
	}
	if !hasAgeErr {
		t.Errorf("expected age range error, got %v", errs)
	}
}

func TestValidateUser_MissingAge(t *testing.T) {
	input := map[string]interface{}{
		"username": "bob",
	}
	errs := validateUser(input)
	hasAgeErr := false
	for _, e := range errs {
		if e.Field == "age" {
			hasAgeErr = true
		}
	}
	if !hasAgeErr {
		t.Errorf("expected age error, got %v", errs)
	}
}

func TestValidateUser_InvalidEmail(t *testing.T) {
	input := map[string]interface{}{
		"username": "charlie",
		"age":      20.0,
		"email":    "not-an-email",
	}
	errs := validateUser(input)
	hasEmailErr := false
	for _, e := range errs {
		if e.Field == "email" {
			hasEmailErr = true
		}
	}
	if !hasEmailErr {
		t.Errorf("expected email error, got %v", errs)
	}
}

func TestValidateUser_InvalidRole(t *testing.T) {
	input := map[string]interface{}{
		"username": "dave",
		"age":      25.0,
		"role":     "superadmin",
	}
	errs := validateUser(input)
	hasRoleErr := false
	for _, e := range errs {
		if e.Field == "role" {
			hasRoleErr = true
		}
	}
	if !hasRoleErr {
		t.Errorf("expected role error, got %v", errs)
	}
}

func TestUserHandler_Success(t *testing.T) {
	body := `{"username":"alice","age":30,"email":"a@b.com","role":"admin"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	userHandler(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestUserHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	userHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestUserHandler_ValidationError(t *testing.T) {
	body := `{"username":"ab","age":10}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	userHandler(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", rec.Code)
	}
	var resp map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&resp)
	if _, ok := resp["errors"]; !ok {
		t.Error("expected errors in response")
	}
}
