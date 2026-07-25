package token

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jm5905938/zjut-work-big/internal/model"
)

const issuer = "zjut-work-big"

type Claims struct {
	UserID   uint64         `json:"user_id"`
	Username string         `json:"username"`
	Role     model.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret    []byte
	expiresIn time.Duration
}

func NewJWTManager(secret string, expiresIn time.Duration) *JWTManager {
	return &JWTManager{
		secret:    []byte(secret),
		expiresIn: expiresIn,
	}
}

func (m *JWTManager) Generate(user *model.User) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(m.expiresIn)

	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   strconv.FormatUint(user.ID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	signedToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	).SignedString(m.secret)
	if err != nil {
		return "", 0, fmt.Errorf("生成 JWT 失败: %w", err)
	}

	return signedToken, int64(m.expiresIn.Seconds()), nil
}

func (m *JWTManager) Parse(signedToken string) (*Claims, error) {
	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(
		signedToken,
		claims,
		func(_ *jwt.Token) (any, error) {
			return m.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("解析 JWT 失败: %w", err)
	}

	if !parsedToken.Valid {
		return nil, errors.New("JWT 无效")
	}

	return claims, nil
}
