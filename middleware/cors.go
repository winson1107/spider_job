package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "content-type", "token"},
		ExposeHeaders:    []string{"Content-Length", "content-type", "token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
