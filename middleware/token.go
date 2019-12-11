package middleware

import (
	"app/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Token() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.DefaultQuery("token", "")
		if len(token) == 0 {
			token = ctx.GetHeader("token")
		}
		user, err := lib.AESDecrypt(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"msg":  "token错误",
				"code": "500",
			})
			return
		}
		ctx.Set("user", string(user))
		ctx.Next()
	}
}
