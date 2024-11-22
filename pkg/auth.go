package pkg

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	jwtSecret []byte
	once      sync.Once
)

func getJWTSecret() []byte {
	once.Do(func() {
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
		log.Println("JWT secret loaded from environment")
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
	log.Println("Hashing password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
	}
	return string(hashedPassword), err
}

func CheckPassword(hashedPassword, password string) error {
	log.Println("Checking password")
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Printf("Password check failed: %v", err)
	}
	return err
}

func GenerateJWT(user User) (string, error) {
	log.Println("Generating JWT for user:", user.Username)
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
	signedToken, err := token.SignedString(getJWTSecret())
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
	}
	return signedToken, err
}

func ParseJWT(tokenStr string) (*Claims, error) {
	log.Println("Parsing JWT")
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		log.Printf("Invalid token: %v", err)
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
