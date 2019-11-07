package handler

import (
	"app/service"
	"github.com/gin-gonic/gin"
	"strconv"
	"regexp"
	"net/http"
	"encoding/json"
	"app/config"
)

var tbk *service.TbkService

//分类
type cate struct {
	CateId   int    `json:"cate_id"`
	CateName string `json:"cate_name"`
}
//生成淘口令参数
type CreateTKlReq struct {
	Url   string `form:"url" json:"url" binding:"required"`
	Text  string `form:"text" json:"text"`
}
func Init() {
	TbkConf := config.Conf.Tbk
	tbk = service.NewTbkService(TbkConf.AppId, TbkConf.AppSecret, TbkConf.AdZoneId)
}
//搜索
func Search(c *gin.Context) {
	keyword := c.DefaultQuery("keyword", "")
	if keyword == "" {
		c.JSON(200, gin.H{"code": 200, "message": "关键词不能为空", "success": false, "data": []string{}})
		return
	}
	keyword = parseKeyword(keyword)
	pageSize, pageNum := getPagination(c)
	params := make(map[string]string)
	params["q"] = keyword
	params["page_size"] = strconv.Itoa(pageSize)
	params["page_no"] = strconv.Itoa(pageNum)
	params["sort"] = getSort(c.DefaultQuery("sort", "0"))
	params["is_tmall"] = c.DefaultQuery("is_tmall", "false")
	params["has_coupon"] = c.DefaultQuery("has_coupon", "false")
	params["need_free_shipment"] = c.DefaultQuery("need_free_shipment", "false")
	resp, err := tbk.Search(params)
	if err != nil {
		c.JSON(200, gin.H{"code": 101, "message": err.Error(), "success": false})
		return
	}
	c.JSON(200, gin.H{"code": 200, "message": "ok", "data": resp, "success": true})
}
//首页
func TbkIndex(c *gin.Context) {
	params := make(map[string]string)
	pageSize, pageNum := getPagination(c)
	params["page_size"] = strconv.Itoa(pageSize)
	params["page_no"] = strconv.Itoa(pageNum)
	cateId, _ := strconv.Atoi(c.DefaultQuery("cate_id", "1"))
	cateType, _ := strconv.Atoi(c.DefaultQuery("cate_type", "1"))
	params["material_id"] = strconv.Itoa(getCate(cateId, cateType))
	resp, err := tbk.Lists(params)
	if err != nil {
		c.JSON(200, gin.H{"code": 101, "message": err.Error(), "success": false})
		return
	}
	c.JSON(200, gin.H{"code": 200, "message": "ok", "data": resp, "success": true})
}
//生成淘口令
func Pwd(c *gin.Context)  {
	var req CreateTKlReq
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(200, gin.H{"code": 101, "message": err.Error(), "success": false})
		return
	}
	params := map[string]string{"url":req.Url,"text":req.Text}
	pwd,err := tbk.CreateTkl(params)
	if err != nil {
		c.JSON(200, gin.H{"code": 101, "message": "链接不能为空", "success": false})
		return
	}
	resp := map[string]string{"pwd":pwd}
	c.JSON(http.StatusOK,gin.H{"code": 200, "message": "ok", "success": true,"data":resp})
}
//获取分页
func getPagination(c *gin.Context) (pageSize, pageNum int) {
	num := c.DefaultQuery("page_num", "1")
	size := c.DefaultQuery("page_size", "10")
	pageSize, _ = strconv.Atoi(size)
	pageNum, _ = strconv.Atoi(num)
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return
}

//获取排序字段
func getSort(sort string) string {
	sortNum, err := strconv.Atoi(sort)
	if err != nil {
		sortNum = 0
	}
	sortArr := []string{
		"total_sales_des",
		"tk_rate_des",
		"price_des",
	}
	if sortNum <= 0 || sortNum > 3 {
		sortNum = 0
	}
	return sortArr[sortNum]
}

//提取关键字
func parseKeyword(keyword string) string {
	var err error
	regs := []string{
		`\$[a-zA-Z0-9]{5,}\$`,
		`＄[a-zA-Z0-9]{5,}＄`,
		`￠[a-zA-Z0-9]{5,}￠`,
		`￡[a-zA-Z0-9]{5,}￡`,
		`￦[a-zA-Z0-9]{5,}￦`,
		`￥[a-zA-Z0-9]{5,}]￥`,
		`([\x{00A0}-\x{02AF}]){1}[a-zA-Z0-9]{5,15}[\x{00A0}-\x{02AF}]`,//拉丁字符
		`([\x{20A0}-\x{20CF}]){1}[a-zA-Z0-9]{5,15}[\x{20A0}-\x{20CF}]`,//货币
	}
	pregSuccess := false
	for _, v := range regs {
		reg := regexp.MustCompile(v)
		if reg.MatchString(keyword) {
			pregSuccess = true
			break
		}
	}
	if pregSuccess {
		keyword, err = tbk.DecryptTkl(keyword)
		if err != nil {
			return keyword
		}
	}
	return keyword
}

//category映射
func getCate(cateId, cateType int) int {
	defaultId := 13366
	data := `[[{"cate_name":"综合","cate_id":13366},{"cate_name":"女装","cate_id":13367},{"cate_name":"男装","cate_id":13372},{"cate_name":"内衣","cate_id":13373},{"cate_name":"母婴","cate_id":13374},{"cate_name":"食品","cate_id":13375},{"cate_name":"家居家装","cate_id":13368},{"cate_name":"鞋包配饰","cate_id":13370},{"cate_name":"美妆个护","cate_id":13371},{"cate_name":"运动户外","cate_id":13376},{"cate_name":"数码家电","cate_id":13369}],[{"cate_name":"综合","cate_id":3786},{"cate_name":"女装","cate_id":3788},{"cate_name":"男装","cate_id":3790},{"cate_name":"内衣","cate_id":3787},{"cate_name":"母婴","cate_id":3789},{"cate_name":"食品","cate_id":3791},{"cate_name":"家居家装","cate_id":3792},{"cate_name":"鞋包配饰","cate_id":3796},{"cate_name":"美妆个护","cate_id":3794},{"cate_name":"运动户外","cate_id":3795},{"cate_name":"数码家电","cate_id":3793}],[{"cate_name":"综合","cate_id":9660},{"cate_name":"女装","cate_id":9658},{"cate_name":"男装","cate_id":9654},{"cate_name":"内衣","cate_id":9652},{"cate_name":"母婴","cate_id":9650},{"cate_name":"食品","cate_id":9649},{"cate_name":"家装","cate_id":9655},{"cate_name":"鞋包配饰","cate_id":9648},{"cate_name":"美妆个护","cate_id":9653},{"cate_name":"运动户外","cate_id":9651},{"cate_name":"数码家电","cate_id":9656}],[{"cate_name":"综合","cate_id":3756},{"cate_name":"女装","cate_id":3767},{"cate_name":"男装","cate_id":3764},{"cate_name":"内衣","cate_id":3765},{"cate_name":"母婴","cate_id":3760},{"cate_name":"食品","cate_id":3761},{"cate_name":"家装","cate_id":3758},{"cate_name":"鞋包配饰","cate_id":3762},{"cate_name":"美妆个护","cate_id":3763},{"cate_name":"运动户外","cate_id":3766},{"cate_name":"数码家电","cate_id":3759}]]`
	rows := make([][]cate, 4)
	err := json.Unmarshal([]byte(data), &rows)
	if err != nil {
		return defaultId
	}
	if cateType > len(rows) {
		return defaultId
	}
	targetCate := rows[cateType-1]
	for i, v := range targetCate {
		if cateId-1 == i {
			return v.CateId
		}
	}
	return defaultId
}
