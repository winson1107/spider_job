package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"app/models"
	"app/lib"
	"net/url"
)

func PostVisit(ctx *gin.Context)  {
	token := ctx.GetHeader("token")
	user,err := lib.AesDec([]byte(token))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,gin.H{
			"msg":"token错误",
			"code":"500",
		})
		return
	}
	visitor := &models.VisitLog{}

	err = ctx.BindJSON(visitor)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":"500",
			"msg":"解析json出错"+err.Error(),
		});
		return
	}
	rows := make([]models.VisitLog,1)
	u,_ := url.Parse(visitor.Url)
	visitor.Host = u.Host
	visitor.User = user
	rows[0] = *visitor
	err = models.InsertLogs(rows)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":"500",
			"msg":"插入出错"+err.Error(),
		});
		return
	}
	ctx.JSON(http.StatusOK,gin.H{
		"code":"200",
		"msg":"ok",
	})
}
