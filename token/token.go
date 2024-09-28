package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jasurxaydarov/book_shop_api_getway/genproto/user_service"
)

type Claim struct {
	UserId   string
	UserRole string
	jwt.StandardClaims
}

var secretJwtKey = []byte("Shit")

func GenerateJWT(claim user_service.Clamis) (string, error) {

	experationTime := time.Now().Add(1 * time.Hour)

	jwtClaim := Claim{
		UserId:         claim.UserId,
		UserRole:       claim.UserRole,
		StandardClaims: jwt.StandardClaims{ExpiresAt: experationTime.Unix()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)

	tokenString, err := token.SignedString(secretJwtKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseJWT(tokenString string) (*Claim, error) {

	var claim = &Claim{}

	token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (interface{}, error) {

		return secretJwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {

		return nil, fmt.Errorf("Invalid token")
	}

	return claim, nil
}
