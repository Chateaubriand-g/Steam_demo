package middleware

import (
	"net/http"
	"steam-backend/config"
	"steam-backend/models"
	"steam-backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "not offer jwt"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "not offer jwt"))
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString, cfg.JWTSecret)
		if err != nil {
			switch err {
			case utils.ErrToKenExpired:
				c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "token expired"))
			case utils.ErrTokenNotValidYet:
				c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "token not vaild yet"))
			case utils.ErrTokenInvalid:
				c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "token is invaild"))
			}
			c.Abort()
			return
		}

		//c.set()将对象存储到gin的上下文，c.next()让请求流转到后续处理，通过c.get()获取存储的对象
		c.Set("userId", claims.UserID)
		c.Next()
	}
}
