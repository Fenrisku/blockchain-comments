package web

// GetPostListHandler2 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]

import (
	"log"
	"net/http"
	"strconv"

	_ "github.com/blockchain.com/comments/docs" // 千万不要忘了导入把你上一步生成的docs

	"github.com/blockchain.com/comments/service"

	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/gin-gonic/gin"
)

// @Summary 评价数据不分类查询
// @Description tracecom?start=&size=
// @Tags 获取评价信息
// @Accept mpfd
// @Produce json
// @Param start query string true "开始数"
// @Param size query string true "查询数量"
// @Success 200 {string} string "{"msg": "SUCCESS"}"
// @Failure 400 {string} string "{"msg": "FAIL}"
// @Router /tracecom [get]
func GetTraceCom(router *gin.Engine, tracecom []service.Comment) {

	router.Handle("GET", "/tracecom", func(c *gin.Context) {
		s := c.Query("start")
		st, err := strconv.Atoi(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		s = c.Query("size")
		sz, err := strconv.Atoi(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		var result []service.Comment
		if sz == 0 {
			result = tracecom[st:]
		} else {
			result = tracecom[st : sz+st]
		}

		tempMap := make(map[string]interface{})
		tempMap["data"] = result
		c.JSON(http.StatusOK, tempMap)
	})
}

// @Summary 获取每门课程评分
// @Description /count
// @Tags 获取评分
// @Accept mpfd
// @Produce json
// @Success 200 {string} string "{"msg": "SUCCESS"}"
// @Failure 400 {string} string "{"msg": "FAIL}"
// @Router /count [get]
func GetCount(router *gin.Engine, count []service.BlockCount) {
	router.Handle("GET", "/count", func(c *gin.Context) {
		tempMap := make(map[string]interface{})
		tempMap["data"] = count
		c.JSON(http.StatusOK, tempMap)
	})
}

// @Summary 获取总评分
// @Description /count/toal
// @Tags 获取评分
// @Accept mpfd
// @Produce json
// @Success 200 {string} string "{"msg": "SUCCESS"}"
// @Failure 400 {string} string "{"msg": "FAIL}"
// @Router /count/total [get]
func GetTotalCount(router *gin.Engine, total service.BlockCount) {
	router.Handle("GET", "/count/total", func(c *gin.Context) {
		tempMap := make(map[string]interface{})
		tempMap["data"] = total
		c.JSON(http.StatusOK, tempMap)
	})
}

// @Summary 获取课程信息
// @Description /tracecom?cid=&start=&size=
// @Tags 获取课程信息
// @Accept mpfd
// @Produce json
// @Param start query string true "开始数"
// @Param size query string true "查询数量"
// @Success 200 {string} string "{"msg": "SUCCESS"}"
// @Failure 400 {string} string "{"msg": "FAIL}"
// @Router /classinfo [get]
func GetClassInfo(router *gin.Engine, classinfo []service.Class) {
	router.Handle("GET", "/classinfo", func(c *gin.Context) {
		s := c.Query("start")
		st, err := strconv.Atoi(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		s = c.Query("size")
		sz, err := strconv.Atoi(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		var result []service.Class
		if sz == 0 {
			result = classinfo[st:]
		} else {
			result = classinfo[st : sz+st]
		}

		tempMap := make(map[string]interface{})
		tempMap["data"] = result
		c.JSON(http.StatusOK, tempMap)
	})
}

// @Summary 获取课程评价
// @Description /classcom?cid=&start=&size=
// @Tags 获取评价信息
// @Accept mpfd
// @Produce json
// @Param cid query string true "cid"
// @Param start query string true "开始数"
// @Param size query string true "查询数量"
// @Success 200 {string} string "{"msg": "SUCCESS"}"
// @Failure 400 {string} string "{"msg": "FAIL}"
// @Router /classcom [get]
func GetClassComment(router *gin.Engine, classcom map[int][]service.Comment) {
	router.Handle("GET", "/classcom", func(c *gin.Context) {
		s := c.Query("cid")
		cid, _ := strconv.Atoi(s)
		s = c.Query("start")
		st, err := strconv.Atoi(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		s = c.Query("size")
		sz, err := strconv.Atoi(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		var result []service.Comment
		if len(classcom[cid]) == 0 || classcom[cid] == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "data not found"})
		}
		if sz == 0 {
			result = classcom[cid][st:]
		} else {
			result = classcom[cid][st : st+sz]
		}
		tempMap := make(map[string]interface{})
		tempMap["data"] = result
		c.JSON(http.StatusOK, tempMap)
	})
}

func WebStart(total service.BlockCount, count []service.BlockCount, tracecom []service.Comment, classinfo []service.Class, classcom map[int][]service.Comment) {

	router := gin.Default()
	//------------朔源----------------------
	router.Use(Cors())
	GetTraceCom(router, tracecom)

	//------------记数----------------------
	GetCount(router, count)
	GetTotalCount(router, total)

	//-------------获取课程信息----------------
	GetClassInfo(router, classinfo)

	//-------------依据课程cid获取评论信息----------------
	GetClassComment(router, classcom)

	router.GET("/api/*any", gs.WrapHandler(swaggerFiles.Handler))

	router.Run(":8000")

}

//关闭跨域，仅用于跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()
		c.Next()
	}
}
