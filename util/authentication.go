package util

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const SecretKey = "MyCC98MyHome"

type CustomJWTClaim struct {
	UserId int64
	jwt.StandardClaims
}

func GenerateJWT() (tokenString string, expireTime time.Time) {
	maxAge := 60 * 60 * 24
	expireTime = time.Now().Add(time.Duration(maxAge) * time.Second)
	myCustomClaim := &CustomJWTClaim{
		UserId: 6,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "CharlieZhang",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myCustomClaim)
	tokenString, err := token.SignedString([]byte(SecretKey))

	if err != nil {
		panic(err)
	}
	return
}
