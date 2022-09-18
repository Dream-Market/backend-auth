package utils

import (
	"errors"
	"time"

	"github.com/Dream-Market/backend-auth/pkg/models"
	"github.com/golang-jwt/jwt"
)

type JwtWrapper struct {
	SecretKey string
	Issuer    string
}

type jwtClaims struct {
	jwt.StandardClaims
	SessionId int64 `json:"sid"`
}

func (w *JwtWrapper) GenerateToken(session models.Session) (signedToken string, err error) {
	claims := &jwtClaims{
		SessionId: session.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: session.ExpiresAt.Unix(),
			Issuer:    w.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err = token.SignedString([]byte(w.SecretKey))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (w *JwtWrapper) ValidateToken(signedToken string) (claims *jwtClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(w.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*jwtClaims)

	if !ok {
		return nil, errors.New("couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("JWT is expired")
	}

	return claims, nil
}
