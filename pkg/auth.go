package pkg

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
	"sync"
	"time"
)

var (
	jwtSecret []byte
	once      sync.Once
)

func GetJWTSecret() []byte {
	once.Do(func() {
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
		logrus.Println("JWT secret loaded from environment")
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
	logrus.Println("Hashing password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Printf("Error hashing password: %v", err)
	}
	return string(hashedPassword), err
}

func CheckPassword(hashedPassword, password string) error {
	logrus.Println("Checking password")
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		logrus.Printf("Password check failed: %v", err)
	}
	return err
}

func GenerateJWT(user User) (string, error) {
	logrus.Println("Generating JWT for user:", user.Username)
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
	signedToken, err := token.SignedString(GetJWTSecret())
	if err != nil {
		logrus.Printf("Error generating JWT: %v", err)
	}
	return signedToken, err
}

func ParseJWT(tokenStr string) (*Claims, error) {
	logrus.Println("Parsing JWT")
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		logrus.Printf("Invalid token: %v", err)
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
