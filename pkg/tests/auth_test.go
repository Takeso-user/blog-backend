package tests

import (
	"os"
	"testing"

	"github.com/Takeso-user/blog-backend/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetJWTSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	secret := pkg.GetJWTSecret()
	assert.Equal(t, []byte("testsecret"), secret)
}

func TestHashPassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := pkg.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
}

func TestCheckPassword(t *testing.T) {
	password := "password123"
	hashedPassword, _ := pkg.HashPassword(password)
	err := pkg.CheckPassword(hashedPassword, password)
	require.NoError(t, err)
}

func TestGenerateJWT(t *testing.T) {
	user := pkg.User{Username: "testuser", Role: "user", Password: "password123"}
	token, err := pkg.GenerateJWT(user)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestParseJWT(t *testing.T) {
	user := pkg.User{Username: "testuser", Role: "user", Password: "password123"}
	token, _ := pkg.GenerateJWT(user)
	claims, err := pkg.ParseJWT(token)
	require.NoError(t, err)
	assert.Equal(t, user.Username, claims.Username)
	assert.Equal(t, user.Role, claims.Role)
}
