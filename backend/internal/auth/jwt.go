package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64 `json:"uid"`
	jwt.RegisteredClaims
}

type JWT struct {
	secret []byte
	exp    time.Duration
}

func New(secret string, exp time.Duration) *JWT {
	return &JWT{secret: []byte(secret), exp: exp}
}

func (j *JWT) Sign(userID int64) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.exp)),
			NotBefore: jwt.NewNumericDate(now.Add(-1 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWT) Parse(tokenStr string) (Claims, error) {
	tk, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})
	if err != nil {
		return Claims{}, err
	}
	claims, ok := tk.Claims.(*Claims)
	if !ok || !tk.Valid {
		return Claims{}, errors.New("invalid token")
	}
	return *claims, nil
}
