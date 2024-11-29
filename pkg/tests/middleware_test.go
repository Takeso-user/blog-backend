package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Takeso-user/blog-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MissingTokenReturnsUnauthorized(t *testing.T) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(pkg.JWTMiddleware())

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(ctx, "GET", "/protected", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing token")
}

func Test_InvalidTokenReturnsUnauthorized(t *testing.T) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(pkg.JWTMiddleware())

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(ctx, "GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

func Test_ValidTokenGrantsAccess(t *testing.T) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(pkg.JWTMiddleware())

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
	})

	token, err := pkg.GenerateJWT(pkg.User{Username: "testuser", Role: "user", Password: "password123"})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(ctx, "GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Access granted")
}
