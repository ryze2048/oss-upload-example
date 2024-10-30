package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"oss-upload-example/middleware"
	"oss-upload-example/router"
	"time"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middleware.Cors())
	var uploadingRouter = router.RouterGroupApp.Uploading

	PublicGroup := Router.Group("")
	{
		// 健康监测
		PublicGroup.GET("/v1/health", func(c *gin.Context) {
			c.Header("Content-Type", "text/event-stream")
			c.Header("Cache-Control", "no-cache")
			c.Header("Connection", "keep-alive")
			// 开始响应数字流
			for i := 1; i <= 100; i++ {
				// 每隔1秒钟发送一个数字
				time.Sleep(time.Millisecond * 50)
				_, _ = fmt.Fprintf(c.Writer, `%d`, i)
				c.Writer.Flush()
			}

			// 响应数字流结束
			c.Status(http.StatusOK)
		})
	}
	{
		uploadingRouter.InitUploadingRouter(PublicGroup)
	}

	return Router
}
