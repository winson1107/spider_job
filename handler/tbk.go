package handler

import (
	"app/service"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"strconv"
)

var tbk *service.TbkService

func init() {
	tbk = service.NewTbkService("28065601", "cfbc2253793f178d6b91631d48e7310e", "109632800369")
}
func TbkIndex(c *gin.Context) {
	keyword := c.DefaultQuery("keyword", "")
	if keyword == "" {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "关键词不能为空",
			"success": false,
			"data":    []interface{}{},
		})
	}
	params := make(map[string]string)
	pageNum := c.GetInt("page_size")
	pageSize := c.GetInt("page_num")
	sort := c.GetInt("sort")
	sortMap := []string{
		"total_sales_des",
		"tk_rate_des",
		"price_des",
	}
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if sort <= 0 {
		sort = 0
	}
	params["sort"] = sortMap[sort]
	params["page_size"] = strconv.Itoa(pageSize)
	params["page_num"] = strconv.Itoa(pageNum)
	params["is_tmall"] = c.DefaultQuery("is_tmall", "false")
	params["has_coupon"] = c.DefaultQuery("has_coupon", "false")
	params["need_free_shipment"] = c.DefaultQuery("need_free_shipment", "false")
	params["q"] = keyword
	resp, err := tbk.Search(params)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    "101",
			"message": err.Error(),
			"success": false,
		})
		return
	}
	for i, v := range resp {
		resp[i].ClickUrl = v.Url
		money, _ := decimal.NewFromString(v.ZkFinalPrice)
		amount, _ := decimal.NewFromString(v.CouponAmount)
		rate, _ := decimal.NewFromString(v.CommissionRate)
		d, _ := decimal.NewFromString("1000000")
		t, _ := decimal.NewFromString("90")
		r := money.Sub(amount).Mul(rate).Mul(t).DivRound(d, 2)
		m, _ := r.Float64()
		resp[i].Commission = m
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "ok",
		"data":    resp,
		"success": true,
	})
}
