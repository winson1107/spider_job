package handler

import (
	"app/models"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type ResponseItem struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Url   string `json:"url"`
	Title string `json:"title"`
	Id    string `json:"id"`
}

type PageParams struct {
	Page     int   ` json:"pageNo" form:"pageNo"`
	PageSize int   ` json:"pageSize" form:"pageSize"`
	Total    int64 `json:"total"`
}

func PostVisit(ctx *gin.Context) {
	user, exist := ctx.Get("user")
	if !exist {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg":  "token错误",
			"code": "500",
		})
		return
	}
	visitor := &models.VisitLog{}
	err := ctx.ShouldBind(visitor)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "解析出错" + err.Error(),
		})
		return
	}
	rows := make([]models.VisitLog, 1)
	u, _ := url.Parse(visitor.Url)
	visitor.Host = u.Host
	visitor.User = user.(string)
	rows[0] = *visitor
	err = models.InsertLogs(rows)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "插入出错" + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "ok",
	})
}

func GetVisit(ctx *gin.Context) {
	params, err := getParams(ctx)
	user, exist := ctx.Get("user")
	if err != nil || !exist {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg":  "参数错误:",
			"code": "500",
		})
		return
	}
	result := models.GetLogs(user.(string), int64(params.Page), int64(params.PageSize))
	items := make([]ResponseItem, 0)
	for _, item := range result {
		items = append(items, ResponseItem{
			Url:   item.Url,
			Start: time.Unix(item.BeginTime, 1).Format("2006-01-02 15:04:05"),
			End:   time.Unix(item.EndTime, 1).Format("2006-01-02 15:04:05"),
			Title: item.Title,
			Id:    item.ID.Hex(),
		})
	}
	//pagination := []int{}
	params.Total = models.GetUserVisitCount(user.(string))
	ctx.JSON(http.StatusOK, gin.H{
		"code":       "200",
		"msg":        "ok",
		"data":       items,
		"pagination": params,
	})
}
