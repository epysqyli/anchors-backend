package middleware

import (
	"net/http"
	"strings"

	"github.com/epysqyli/anchors-backend/domain"
	"github.com/epysqyli/anchors-backend/internal/tokenutil"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// keep a way for testing purpose to use curl + Auth header?
		// authHeader := c.Request.Header.Get("Authorization")
		// t := strings.Split(authHeader, " ")

		tokenCookie, err := c.Cookie(domain.AuthToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
			c.Abort()
			return
		}

		accessToken := strings.Split(tokenCookie, "---")[0]
		authorized, err := tokenutil.IsAuthorized(accessToken, secret)

		if authorized {
			userID, err := tokenutil.ExtractIDFromToken(accessToken, secret)

			if err != nil {
				c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
				c.Abort()
				return
			}

			c.Set("x-user-id", userID)
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
		c.Abort()
		return
	}
}
