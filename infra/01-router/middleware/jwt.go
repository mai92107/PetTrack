// middleware/jwt.go
package middleware

import (
	"PetTrack/infra/00-core/util/logafa"
	"PetTrack/infra/00-core/model"
	jwtUtil "PetTrack/infra/00-core/util/jwt"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Permission int

const (
	PermGuest Permission = iota
	PermMember
	PermAdmin
)

var (
	ErrMissingToken = errors.New("missing token")
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
	ErrForbidden    = errors.New("permission denied")
)

func JWTMiddleware(required Permission) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 取得 Token
		token := c.GetHeader("jwt")
		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = token[7:]
		}

		// 驗證 token
		claims, err := ValidateJWT(token, required)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 成功，把 claims 設進 context
		c.Set("claims", claims)
		c.Next()
	}
}

func ValidateJWT(jwt string, required Permission) (*model.Claims, error) {

	if jwt == "" {
		return nil, ErrMissingToken
	}

	// 1. Parse token
	claims, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Warn("JWT 解析失敗", "error", err, "token", MaskToken(jwt))
		return nil, ErrInvalidToken
	}

	// 2. expiration
	if !claims.ExpiresAt.After(time.Now().UTC()) {
		logafa.Warn("JWT 過期", "token", MaskToken(jwt))
		return nil, ErrExpiredToken
	}

	// 3. permission check
	switch required {
	case PermAdmin:
		if !claims.IsAdmin() {
			logafa.Warn("權限不足", "member", claims.MemberId)
			return nil, ErrForbidden
		}
	case PermMember:
		// guest 無法登入，member 與 admin 通過
		// 已驗證 token，所以 Member/Admin 都允許
	}

	return claims, nil
}

// 遮蔽 token（安全日誌）
func MaskToken(token string) string {
	if len(token) < 10 {
		return "****"
	}
	return token[:6] + "..." + token[len(token)-4:]
}
