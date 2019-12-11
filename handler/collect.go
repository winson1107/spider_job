package handler

import (
	"app/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Time time.Time

const (
	timeFormat = "2006-01-02 15:04:05"
)

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}
func (t Time) String() string {
	return time.Time(t).Format(timeFormat)
}

type QueryParams struct {
	Start time.Time `json:"start" form:"start" binding:"required" time_format:"2006-01-02 15:04:05" time_utc:"1" binding:"required"`
	End   time.Time `json:"end" form:"end" binding:"required" time_format:"2006-01-02 15:04:05" time_utc:"1" binding:"required"`
}

func getParams(ctx *gin.Context) (*PageParams, error) {
	params := &PageParams{}
	err := ctx.ShouldBindQuery(params)
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	return params, err
}
func GetCollectUrl(ctx *gin.Context) {
	pageParams, err := getParams(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"code": "500",
		})
		return
	}
	items := models.GetCollects(int64(pageParams.Page), int64(pageParams.PageSize))
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "ok",
		"data": items,
	})
}

func PostCollectUrl(ctx *gin.Context) {
	item := &models.CollectItem{}
	err := ctx.ShouldBindJSON(item)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"code": "500",
		})
		return
	}
	//忽略错误
	models.AddCollectUrl(item)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "ok",
	})
}

func QueryCollect(ctx *gin.Context) {
	query := &QueryParams{}
	err := ctx.ShouldBindQuery(query)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"code": "500",
		})
		return
	}
	items := models.QueryCollect(query.Start, query.End)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "ok",
		"data": items,
	})
}
