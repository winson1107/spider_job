package main

import (
	"app/config"
	"app/handler"
	"app/middleware"
	"app/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	APIVersion = "v2"
)

func AddRouter(engine *gin.Engine) {
	engine.Use(middleware.Cors())
	g := engine.Group("/" + APIVersion)
	g.Any("/", handler.Index)
	g.GET("/company", handler.CompanyList)
	g.POST("/company", handler.AddCompany)
	g.POST("/position", handler.AddPosition)
	g.GET("/company/hot", handler.GetHotCompany)
	g.GET("/company/job/:company_id", handler.CompanyJobs)
	g.GET("/company/analysis", handler.CompanyAnaly)
	g.GET("/position/work_year", handler.WorkYear)
	g.GET("/position/edu", handler.Educational)
	g.GET("/position/weekly", handler.Weekly)
	g.GET("/position/monthly", handler.Monthly)
	g.GET("/word", func(c *gin.Context) {
		data, _ := ioutil.ReadFile("./word.json")
		result := make(map[string]interface{}, 0)
		_ = json.Unmarshal(data, &result)
		c.JSON(200, result)
	})
	g.GET("/tbk", handler.TbkIndex)
	g.GET("/tbk/search", handler.Search)
	g.POST("/tbk/pwd", handler.Pwd)
	g.GET("/user/info", func(ctx *gin.Context) {
		js := `{"id":"4291d7da9005377ec9aec4a71ea837f","name":"天野远子","username":"admin","password":"","avatar":"/avatar2.jpg","status":1,"telephone":"","lastLoginIp":"27.154.74.117","lastLoginTime":1534837621348,"creatorId":"admin","createTime":1497160610259,"merchantCode":"TLif2btpzg079h15bk","deleted":0,"roleId":"admin","role":{"id":"admin","name":"管理员","describe":"拥有所有权限","status":1,"creatorId":"system","createTime":1497160610259,"deleted":0,"permissions":[{"roleId":"admin","permissionId":"dashboard","permissionName":"仪表盘"},{"roleId":"admin","permissionId":"table","permissionName":"表格权限"},{"roleId":"admin","permissionId":"form","permissionName":"表单权限"},{"roleId":"admin","permissionId":"user","permissionName":"用户管理"}]}}`
		h := gin.H{}
		json.Unmarshal([]byte(js), &h)
		ctx.JSON(http.StatusOK, h)
	})
	g.Use(middleware.Token()).POST("/visit", handler.PostVisit)
	g.Use(middleware.Token()).GET("/visit", handler.GetVisit)
	g.Use(middleware.Token()).GET("/collect", handler.GetCollectUrl)
	g.Use(middleware.Token()).POST("/collect", handler.PostCollectUrl)
	g.Use(middleware.Token()).GET("/collect/query", handler.QueryCollect)

}

func Run() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	AddRouter(engine)
	log.Fatal(engine.Run(config.Conf.Ip + ":" + strconv.Itoa(config.Conf.Port)))
}
func main() {
	//初始化配置
	err := config.Init("./config.json")
	handler.Init()
	if err != nil {
		log.Fatal(err)
	}
	models.Init()
	Run()
}
