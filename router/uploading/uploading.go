package uploading

import "github.com/gin-gonic/gin"

type UploadingRouter struct{}

func (u *UploadingRouter) InitUploadingRouter(Router *gin.RouterGroup) {
	var v1 = Router.Group("v1")
	uploadingWithoutRecord := v1.Group("uploading")
	{
		uploadingWithoutRecord.POST("oss", func(c *gin.Context) {
			c.JSON(200, gin.H{"state": "ok"})
		})
	}
}
