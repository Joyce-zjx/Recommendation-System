package token

import (
	"errors"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrJWTExpired           = errors.New("jwt expired")
)

var (
	jwtkey        = []byte("JonathanJoyceAreCool")
	signingMethod = jwt.SigningMethodHS256
)

type Claims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenJWT(userID string, username string, expiryTime int64) (string, error) {
	claims := Claims{
		UserID:         userID,
		Username:       username,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expiryTime},
	}
	return jwt.NewWithClaims(signingMethod, claims).SignedString(jwtkey)
}

func ParseJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return jwtkey, nil
	})
	if err != nil {
		// expired token will be identified in this block, then need user to re-login
		if strings.Contains(err.Error(), "expired") {
			return nil, ErrJWTExpired
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func RefreshJWT(tokenStr string, expiryTime int64) (string, error) {
	claims, err := ParseJWT(tokenStr)
	if err != nil {
		return "", err
	}
	claims.ExpiresAt = expiryTime
	return jwt.NewWithClaims(signingMethod, claims).SignedString(jwtkey)
}
