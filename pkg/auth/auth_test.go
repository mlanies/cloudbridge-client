package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateJWT(secret string, claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func TestValidateJWTToken_Valid(t *testing.T) {
	secret := "testsecret"
	claims := jwt.MapClaims{
		"sub": "user1",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	tokenStr, err := generateJWT(secret, claims)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	am, err := NewAuthManager(&AuthConfig{Type: "jwt", Secret: secret})
	if err != nil {
		t.Fatalf("failed to create auth manager: %v", err)
	}
	_, err = am.ValidateToken(tokenStr)
	if err != nil {
		t.Errorf("expected valid token, got error: %v", err)
	}
}

func TestValidateJWTToken_Expired(t *testing.T) {
	secret := "testsecret"
	claims := jwt.MapClaims{
		"sub": "user1",
		"exp": time.Now().Add(-time.Hour).Unix(),
	}
	tokenStr, err := generateJWT(secret, claims)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	am, err := NewAuthManager(&AuthConfig{Type: "jwt", Secret: secret})
	if err != nil {
		t.Fatalf("failed to create auth manager: %v", err)
	}
	_, err = am.ValidateToken(tokenStr)
	if err == nil {
		t.Errorf("expected error for expired token, got nil")
	}
}

func TestValidateJWTToken_InvalidSignature(t *testing.T) {
	secret := "testsecret"
	claims := jwt.MapClaims{
		"sub": "user1",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	tokenStr, err := generateJWT("wrongsecret", claims)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	am, err := NewAuthManager(&AuthConfig{Type: "jwt", Secret: secret})
	if err != nil {
		t.Fatalf("failed to create auth manager: %v", err)
	}
	_, err = am.ValidateToken(tokenStr)
	if err == nil {
		t.Errorf("expected error for invalid signature, got nil")
	}
}

func TestValidateJWTToken_NoSubClaim(t *testing.T) {
	secret := "testsecret"
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	tokenStr, err := generateJWT(secret, claims)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	am, err := NewAuthManager(&AuthConfig{Type: "jwt", Secret: secret})
	if err != nil {
		t.Fatalf("failed to create auth manager: %v", err)
	}
	tok, err := am.ValidateToken(tokenStr)
	if err != nil {
		t.Errorf("expected valid token, got error: %v", err)
	}
	_, err = am.ExtractSubject(tok)
	if err == nil {
		t.Errorf("expected error for missing sub claim, got nil")
	}
}

func TestExtractTenantID(t *testing.T) {
	// Create a test token with tenant_id
	claims := jwt.MapClaims{
		"sub":       "test-user",
		"tenant_id": "tenant-001",
		"exp":       time.Now().Add(time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := &jwt.Token{
		Claims: claims,
		Valid:  true,
	}

	am := &AuthManager{}

	// Test extraction
	tenantID, err := am.ExtractTenantID(token)
	if err != nil {
		t.Fatalf("Failed to extract tenant_id: %v", err)
	}

	if tenantID != "tenant-001" {
		t.Errorf("Expected tenant_id 'tenant-001', got '%s'", tenantID)
	}
}

func TestExtractTenantIDMissing(t *testing.T) {
	// Create a test token without tenant_id
	claims := jwt.MapClaims{
		"sub": "test-user",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := &jwt.Token{
		Claims: claims,
		Valid:  true,
	}

	am := &AuthManager{}

	// Test extraction of missing tenant_id
	tenantID, err := am.ExtractTenantID(token)
	if err != nil {
		t.Fatalf("Failed to extract tenant_id: %v", err)
	}

	if tenantID != "" {
		t.Errorf("Expected empty tenant_id, got '%s'", tenantID)
	}
}

func TestExtractClaims(t *testing.T) {
	// Create a test token with both subject and tenant_id
	claims := jwt.MapClaims{
		"sub":       "test-user",
		"tenant_id": "tenant-001",
		"exp":       time.Now().Add(time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := &jwt.Token{
		Claims: claims,
		Valid:  true,
	}

	am := &AuthManager{}

	// Test extraction of both claims
	subject, tenantID, err := am.ExtractClaims(token)
	if err != nil {
		t.Fatalf("Failed to extract claims: %v", err)
	}

	if subject != "test-user" {
		t.Errorf("Expected subject 'test-user', got '%s'", subject)
	}

	if tenantID != "tenant-001" {
		t.Errorf("Expected tenant_id 'tenant-001', got '%s'", tenantID)
	}
}
