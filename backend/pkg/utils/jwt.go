package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

// Claims JWT claims
type Claims struct {
	UserID   uint64 `json:"userId,string"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// MarshalJSON 将雪花 ID 以字符串形式写入 JWT 载荷，避免前端解析时丢失精度。
func (c Claims) MarshalJSON() ([]byte, error) {
	type claimsJSON struct {
		UserID   string `json:"userId"`
		Username string `json:"username"`
		jwt.RegisteredClaims
	}

	return json.Marshal(claimsJSON{
		UserID:           strconv.FormatUint(c.UserID, 10),
		Username:         c.Username,
		RegisteredClaims: c.RegisteredClaims,
	})
}

// UnmarshalJSON 兼容字符串和历史数字形式的 JWT 用户 ID。
func (c *Claims) UnmarshalJSON(data []byte) error {
	type claimsJSON struct {
		UserID   json.RawMessage `json:"userId"`
		Username string          `json:"username"`
		jwt.RegisteredClaims
	}

	var raw claimsJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	userID, err := parseUint64JSON(raw.UserID)
	if err != nil {
		return fmt.Errorf("invalid userId: %w", err)
	}

	c.UserID = userID
	c.Username = raw.Username
	c.RegisteredClaims = raw.RegisteredClaims
	return nil
}

func parseUint64JSON(data json.RawMessage) (uint64, error) {
	if len(data) == 0 {
		return 0, nil
	}

	var stringID string
	if err := json.Unmarshal(data, &stringID); err == nil {
		return strconv.ParseUint(stringID, 10, 64)
	}

	var numericID uint64
	if err := json.Unmarshal(data, &numericID); err != nil {
		return 0, err
	}
	return numericID, nil
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint64, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// SetJWTSecret 设置JWT密钥（用于测试）
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}
