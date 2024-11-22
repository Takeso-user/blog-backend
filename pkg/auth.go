package pkg

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"

	"time"
)

var (
	jwtSecret []byte
	once      sync.Once
)

func getJWTSecret() []byte {
	once.Do(func() {
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	})
	return jwtSecret
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateJWT(user User) (string, error) {

	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		Password: user.Password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func ParseJWT(tokenStr string) (*Claims, error) {

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
