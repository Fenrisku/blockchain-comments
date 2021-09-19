package web

import (
	"log"
	"net/http"
	"strconv"

	"github.com/blockchain.com/comments/service"

	"github.com/gin-gonic/gin"
)

// 启动Web服务并指定路由信息
func WebStart(total service.BlockCount, count []service.BlockCount, tracecom []service.Comment, classinfo []service.Class, classcom map[int][]service.Comment) {

	router := gin.Default()

	//------------朔源----------------------
	router.Use(Cors())

	//tracecom?start=&size=
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

	//------------记数----------------------

	router.Handle("GET", "/count", func(c *gin.Context) {
		tempMap := make(map[string]interface{})
		tempMap["data"] = count
		c.JSON(http.StatusOK, tempMap)
	})
	router.Handle("GET", "/count/total", func(c *gin.Context) {
		tempMap := make(map[string]interface{})
		tempMap["data"] = total
		c.JSON(http.StatusOK, tempMap)
	})

	//-------------商店信息----------------
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

	//-------------评论信息----------------
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

	router.Run(":8000")

}

//关闭跨域，仅用于测试
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
