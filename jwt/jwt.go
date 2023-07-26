package jwt

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type CustomJWTClaims struct {
	UserId     int    `json:"userId"`
	LoginId    string `json:"loginId"`
	Username   string `json:"username"`
	ClientId   string `json:"clientId"`
	ProfileKey string `json:"profileKey"`
	jwt.RegisteredClaims
}

func GenJWT(userId int, loginId, userName, clientId, signKey, profileKey, issuer string, duration time.Duration) (string, error) {
	now := time.Now()
	exp := now.Add(duration)

	_uuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	claims := CustomJWTClaims{
		UserId:     userId,
		LoginId:    loginId,
		Username:   userName,
		ClientId:   clientId,
		ProfileKey: profileKey,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        _uuid.String(),
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    issuer,
			Audience: jwt.ClaimStrings{
				clientId,
			},
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(signKey))
}

func ParseJWT(tokenString, signKey string) (*CustomJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomJWTClaims{}, func(token *jwt.Token) (i any, e error) {
		return []byte(signKey), nil
	})
	if err != nil {
		return nil, err
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomJWTClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, nil

	} else {
		return nil, nil
	}
}
