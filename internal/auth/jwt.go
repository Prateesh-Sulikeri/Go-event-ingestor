package auth

import (
	"time"

	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret []byte
	issuer string
	exp    time.Duration
}

func NewJWTService(cfg config.Config) *JWTService {
	return &JWTService{
		secret: []byte(cfg.JWT_SECRET),
		issuer: cfg.JWT_ISSUER,
		exp:    time.Duration(cfg.JWT_EXP_HOURS) * time.Hour,
	}
}

func (j *JWTService) GenerateToken(clientID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": clientID,          // subject of token (our client identifier)
		"iss": j.issuer,          // issuer
		"exp": time.Now().Add(j.exp).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}
