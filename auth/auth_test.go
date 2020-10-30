package auth

import (
	"testing"
)

const (
	userID       = "b7d718aa-1c3f-4367-a35b-bbf951ab8e13"
	tokenVersion = 5
)

func TestGenerateAccessToken(t *testing.T) {
	tokenStr, err := GenerateAccessToken(userID)
	if err != nil {
		t.Error("failed to generate token:", err)
	}

	claims := AccessTokenClaims{}
	token, err := ParseToken(tokenStr, &claims)
	if err != nil {
		t.Fatal("failed to parse token:", err)
	}

	if !token.Valid {
		t.Error("token is not valid")
	}

	if claims.UserID != userID {
		t.Errorf("expected claims.UserID to be \"%s\" but got \"%s\"", userID, claims.UserID)
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	tokenStr, err := GenerateRefreshToken(userID, 5)
	if err != nil {
		t.Error("failed to generate token:", err)
	}

	claims := RefreshTokenClaims{}
	token, err := ParseToken(tokenStr, &claims)
	if err != nil {
		t.Fatal("failed to parse token:", err)
	}

	if !token.Valid {
		t.Error("token is not valid")
	}

	if claims.UserID != userID {
		t.Errorf("expected claims.UserID to be \"%s\" but got \"%s\"", userID, claims.UserID)
	}

	if claims.TokenVersion != tokenVersion {
		t.Errorf("expected claims.TokenVersion to be '%d' but got '%d'", tokenVersion, claims.TokenVersion)
	}
}
