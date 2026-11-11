package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateAPIKey_ValidKey(t *testing.T) {
	key := strings.Repeat("a", 40)
	if err := ValidateAPIKey(key); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestValidateAPIKey_TooShort(t *testing.T) {
	key := "short"
	if err := ValidateAPIKey(key); err == nil {
		t.Error("expected error for short key")
	} else if !strings.Contains(err.Error(), "too short") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateAPIKey_InvalidCharacter(t *testing.T) {
	key := strings.Repeat("a", 35) + "!"
	if err := ValidateAPIKey(key); err == nil {
		t.Error("expected error for invalid character")
	} else if !strings.Contains(err.Error(), "invalid character") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateAPIKey_Compromised(t *testing.T) {
	key := "compromised-key-1234567890abcdef1234567890"
	if err := ValidateAPIKey(key); err == nil {
		t.Error("expected error for compromised key")
	} else if !strings.Contains(err.Error(), "compromised") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateAPIKey_Empty(t *testing.T) {
	if err := ValidateAPIKey(""); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestValidateCreateUser_Valid(t *testing.T) {
	req := CreateUserRequest{
		Username: "alice123",
		Email:    "alice@example.com",
		Password: "secret123",
	}
	errs := validateCreateUser(req)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateCreateUser_InvalidEmail(t *testing.T) {
	req := CreateUserRequest{
		Username: "alice",
		Email:    "not-an-email",
		Password: "secret123",
	}
	errs := validateCreateUser(req)
	if len(errs) == 0 {
		t.Fatal("expected errors")
	}
	found := false
	for _, e := range errs {
		if e.Field == "email" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected email error, got %v", errs)
	}
}

func TestValidateCreateUser_ShortPassword(t *testing.T) {
	req := CreateUserRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "short",
	}
	errs := validateCreateUser(req)
	found := false
	for _, e := range errs {
		if e.Field == "password" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected password error, got %v", errs)
	}
}

func TestValidateName_Empty(t *testing.T) {
	if err := validateName(""); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestValidateName_InvalidChars(t *testing.T) {
	if err := validateName("<script>"); err == nil {
		t.Error("expected error for invalid name")
	}
}

func TestValidateName_Valid(t *testing.T) {
	if err := validateName("alice_123"); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestCreateUserHandler_JSONResponses(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantCT         string
		wantErrorSub   string
		wantCreatedUsr string
	}{
		{
			name:         "unknown field produces valid JSON error",
			method:       http.MethodPost,
			body:         `{"username":"alice","email":"alice@example.com","password":"secret123","role":"admin"}`,
			wantStatus:   http.StatusBadRequest,
			wantCT:       "application/json",
			wantErrorSub: `unknown field "role"`,
		},
		{
			name:         "malformed JSON produces valid JSON error",
			method:       http.MethodPost,
			body:         `{"username":`,
			wantStatus:   http.StatusBadRequest,
			wantCT:       "application/json",
			wantErrorSub: "invalid JSON:",
		},
		{
			name:           "valid payload returns created JSON",
			method:         http.MethodPost,
			body:           `{"username":"alice","email":"alice@example.com","password":"secret123"}`,
			wantStatus:     http.StatusCreated,
			wantCT:         "application/json",
			wantCreatedUsr: "alice",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/users", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			createUserHandler(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if gotCT := rec.Header().Get("Content-Type"); !strings.Contains(gotCT, tc.wantCT) {
				t.Errorf("Content-Type = %q, want %q", gotCT, tc.wantCT)
			}

			var payload map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
				t.Fatalf("response is not valid JSON (%v): %s", err, rec.Body.String())
			}
			if tc.wantErrorSub != "" {
				errMsg, _ := payload["error"].(string)
				if !strings.Contains(errMsg, tc.wantErrorSub) {
					t.Errorf("error = %q, want substring %q", errMsg, tc.wantErrorSub)
				}
			}
			if tc.wantCreatedUsr != "" && payload["username"] != tc.wantCreatedUsr {
				t.Errorf("username = %v, want %q", payload["username"], tc.wantCreatedUsr)
			}
		})
	}
}
