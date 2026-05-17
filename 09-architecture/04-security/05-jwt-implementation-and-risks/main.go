// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: JWT - implementation and risks
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn what a JWT contains, how signing works, and why tokens still create real operational risk when used carelessly.
//
// WHY THIS MATTERS:
//   - A JWT is a signed claim set, not a trust system by itself.
//
// RUN:
//   go run ./09-architecture/04-security/05-jwt-implementation-and-risks
//
// KEY TAKEAWAY:
//   - Learn what a JWT contains, how signing works, and why tokens still create real operational risk when used carelessly.
// ============================================================================

package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Token (Struct): holds the three JWT segments after parsing — header, claims, and signature
//
// NOTE: This is a teaching implementation. Production systems should use a
// mature library to avoid subtle security bugs like timing leaks or algorithm
// confusion.
type Token struct {
	Header    map[string]any
	Claims    map[string]any
	Signature []byte
	Raw       string
}

// secret (Variable): simulates a server-side HMAC signing key
//
// ROLE: In production, load this from a secure vault or env variable, never
// bake it into source code.
var secret = []byte("super-secret-key-with-at-least-32-characters!!")

// encodeSegment (Function): base64url-encodes a JSON-marshalled value
func encodeSegment(v any) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

// decodeSegment (Function): base64url-decodes and JSON-unmarshals a JWT segment
func decodeSegment(seg string, v any) error {
	data, err := base64.RawURLEncoding.DecodeString(seg)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// sign (Function): computes an HMAC-SHA256 signature for the given signing input
//
// BOUNDARY: This is the core integrity primitive. Without this, anyone can forge claims.
func sign(input string, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(input))
	return mac.Sum(nil)
}

// createToken (Function): builds and signs a complete JWT string
//
// INVARIANT: The returned token is in standard three-segment format:
// <base64url(header)>.<base64url(claims)>.<base64url(signature)>
func createToken(header, claims map[string]any) (string, error) {
	headerSeg, err := encodeSegment(header)
	if err != nil {
		return "", fmt.Errorf("encode header: %w", err)
	}
	claimsSeg, err := encodeSegment(claims)
	if err != nil {
		return "", fmt.Errorf("encode claims: %w", err)
	}
	signingInput := headerSeg + "." + claimsSeg
	sig := sign(signingInput, secret)
	return signingInput + "." + base64.RawURLEncoding.EncodeToString(sig), nil
}

// parseUnsafe (Function): decodes a JWT without any signature verification
//
// WARNING: This is intentionally vulnerable. It shows why you MUST always
// verify the signature before trusting claims.
func parseUnsafe(tokenStr string) (*Token, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("token must have 3 parts, got %d", len(parts))
	}

	var header map[string]any
	if err := decodeSegment(parts[0], &header); err != nil {
		return nil, fmt.Errorf("decode header: %w", err)
	}

	var claims map[string]any
	if err := decodeSegment(parts[1], &claims); err != nil {
		return nil, fmt.Errorf("decode claims: %w", err)
	}

	sig, _ := base64.RawURLEncoding.DecodeString(parts[2])

	return &Token{
		Header:    header,
		Claims:    claims,
		Signature: sig,
		Raw:       tokenStr,
	}, nil
}

// verify (Function): validates the HMAC-SHA256 signature and enforces algorithm policy
//
// INVARIANT: Regects tokens whose "alg" is not "HS256", including the infamous
// "none" algorithm that attackers use to bypass verification.
func verify(tokenStr string) (*Token, error) {
	tok, err := parseUnsafe(tokenStr)
	if err != nil {
		return nil, err
	}

	alg, _ := tok.Header["alg"].(string)
	if alg != "HS256" {
		return nil, fmt.Errorf("unexpected signing algorithm %q — only HS256 is accepted", alg)
	}

	parts := strings.SplitN(tokenStr, ".", 3)
	signingInput := parts[0] + "." + parts[1]

	expectedSig := sign(signingInput, secret)

	gotSig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("decode signature: %w", err)
	}

	if !hmac.Equal(gotSig, expectedSig) {
		return nil, fmt.Errorf("signature mismatch — token has been tampered with")
	}

	return tok, nil
}

// vulnerableVerify (Function): demonstrates the classic "alg: none" JWT attack
//
// FAILURE MODE: This function trusts whatever algorithm the header declares.
// An attacker changes "alg":"HS256" to "alg":"none", removes the signature,
// and this function silently accepts forged claims.
func vulnerableVerify(tokenStr string) (*Token, error) {
	tok, err := parseUnsafe(tokenStr)
	if err != nil {
		return nil, err
	}

	alg, _ := tok.Header["alg"].(string)

	// VULNERABILITY: An attacker sets alg to "none" to bypass all verification.
	if strings.EqualFold(alg, "none") {
		fmt.Println("  [!] VULNERABLE: accepted token with alg='none' — no signature checked!")
		return tok, nil
	}

	return verify(tokenStr)
}

// printToken (Function): pretty-prints a token's contents
func printToken(tok *Token) {
	headerJSON, _ := json.MarshalIndent(tok.Header, "    ", "  ")
	claimsJSON, _ := json.MarshalIndent(tok.Claims, "    ", "  ")
	fmt.Printf("    Header:   %s\n", string(headerJSON))
	fmt.Printf("    Claims:   %s\n", string(claimsJSON))
	fmt.Printf("    Raw:      %s\n", tok.Raw)
}

func main() {
	fmt.Println("=== SEC.5 JWT - implementation and risks ===")
	fmt.Println()

	// ---- 1. Create a signed token ----
	fmt.Println("--- 1. Creating a signed token ---")

	header := map[string]any{
		"alg": "HS256",
		"typ": "JWT",
	}
	claims := map[string]any{
		"sub":   "user_42",
		"name":  "Alice",
		"role":  "admin",
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour).Unix(),
		"admin": true,
	}

	tokenStr, err := createToken(header, claims)
	if err != nil {
		fmt.Printf("  [ERROR] createToken: %v\n", err)
		return
	}
	fmt.Println("  Generated token:")
	fmt.Printf("    %s\n", tokenStr)
	fmt.Println()

	// ---- 2. Verify a valid token ----
	fmt.Println("--- 2. Verifying a valid token ---")

	parsed, err := verify(tokenStr)
	if err != nil {
		fmt.Printf("  [ERROR] verify: %v\n", err)
	} else {
		fmt.Println("  [OK] Token verified successfully!")
		printToken(parsed)
	}
	fmt.Println()

	// ---- 3. Reject a tampered token ----
	fmt.Println("--- 3. Rejecting a tampered token ---")

	parts := strings.Split(tokenStr, ".")
	if len(parts) == 3 {
		// Decode claims, flip admin to false, re-encode
		var tamperedClaims map[string]any
		if err := decodeSegment(parts[1], &tamperedClaims); err == nil {
			tamperedClaims["role"] = "user"
			tamperedClaims["admin"] = false
			newClaimsSeg, err := encodeSegment(tamperedClaims)
			if err == nil {
				tamperedToken := parts[0] + "." + newClaimsSeg + "." + parts[2]
				_, err = verify(tamperedToken)
				if err != nil {
					fmt.Printf("  [OK] Tampered token correctly rejected: %v\n", err)
				}
			}
		}
	}
	fmt.Println()

	// ---- 4. Demonstrate the "alg: none" attack ----
	fmt.Println("--- 4. The 'alg: none' attack ---")
	fmt.Println("  Attacker creates a token with alg='none' and arbitrary claims:")

	attackHeader := map[string]any{
		"alg": "none",
		"typ": "JWT",
	}
	attackClaims := map[string]any{
		"sub":   "attacker",
		"name":  "Eve",
		"role":  "admin",
		"admin": true,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	// For the none attack, the token has no signature (empty signature part)
	headerSeg, _ := encodeSegment(attackHeader)
	claimsSeg, _ := encodeSegment(attackClaims)
	attackToken := headerSeg + "." + claimsSeg + "."

	fmt.Printf("  Attack token: %s\n", attackToken)
	fmt.Println()

	// Safe verification rejects it
	fmt.Println("  Safe verify() result:")
	_, err = verify(attackToken)
	if err != nil {
		fmt.Printf("    [OK] Rejected: %v\n", err)
	} else {
		fmt.Println("    [PANIC] Token accepted (this should never happen!)")
	}
	fmt.Println()

	// Vulnerable verification accepts it
	fmt.Println("  Vulnerable verify() result:")
	vulnTok, err := vulnerableVerify(attackToken)
	if err != nil {
		fmt.Printf("    [OK] Rejected: %v\n", err)
	} else {
		fmt.Println("    [VULNERABLE] Token accepted!")
		printToken(vulnTok)
	}
	fmt.Println()

	// ---- 5. Why you must always verify ----
	fmt.Println("--- 5. parseUnsafe (no verification) ---")
	fmt.Println("  Decoding a forged token without checking signature:")

	forgedClaims := map[string]any{
		"sub":   "impostor",
		"name":  "Mallory",
		"role":  "superadmin",
		"admin": true,
		"exp":   time.Now().Add(365 * 24 * time.Hour).Unix(),
	}
	forgedSeg, _ := encodeSegment(forgedClaims)
	forgedToken := headerSeg + "." + forgedSeg + "." // alg=none, no signature

	unsafeTok, _ := parseUnsafe(forgedToken)
	if unsafeTok != nil {
		fmt.Println("  [WARNING] parseUnsafe decoded the token without verifying:")
		printToken(unsafeTok)
		fmt.Println()
		fmt.Println("  Moral: always verify the signature before acting on claims.")
	}
	fmt.Println()

	fmt.Println("- Signing proves integrity, not that every claim is safe to trust blindly.")
	fmt.Println("- Validate issuer, audience, expiry, and algorithm policy.")
	fmt.Println("- Treat tokens as credentials and logs as hostile to secret material.")
	fmt.Println()
	fmt.Println("JWT mistakes are rarely library mistakes. They are usually trust-boundary mistakes like weak key policy, bad expiry handling, or missing audience checks.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SEC.6 -> 09-architecture/04-security/06-password-hashing")
	fmt.Println("Current: SEC.5 (jwt - implementation and risks)")
	fmt.Println("---------------------------------------------------")
}
