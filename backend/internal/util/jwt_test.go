package util

import (
	"testing"
	"time"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret"
	token, err := GenerateToken(1, "alice", "user", secret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if claims.UserID != 1 || claims.Username != "alice" || claims.Role != "user" {
		t.Errorf("unexpected claims: %+v", claims)
	}
}

func TestParseTokenWrongSecret(t *testing.T) {
	token, _ := GenerateToken(1, "alice", "user", "a", time.Hour)
	if _, err := ParseToken(token, "b"); err == nil {
		t.Error("expected error for wrong secret")
	}
}

func TestParseTokenExpired(t *testing.T) {
	token, _ := GenerateToken(1, "alice", "user", "s", -time.Minute)
	_, err := ParseToken(token, "s")
	if err == nil {
		t.Error("expected expired token error")
	}
}
