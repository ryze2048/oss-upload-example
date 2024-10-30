package uploading

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oss-upload-example/model/common/response"
	"oss-upload-example/utils/upload"
)

type UploadingRouter struct{}

func (u *UploadingRouter) InitUploadingRouter(Router *gin.RouterGroup) {
	var v1 = Router.Group("v1")
	uploadingWithoutRecord := v1.Group("uploading")
	{
		uploadingWithoutRecord.GET("oss", func(c *gin.Context) {
			var err error
			var oss = upload.NewOss()
			var result *upload.PolicyToken
			if result, err = oss.GetPolicyToken(); err != nil {
				c.JSON(http.StatusBadRequest, "")
				return
			}
			c.JSON(http.StatusOK, result)
		})

		uploadingWithoutRecord.GET("callback", func(c *gin.Context) {
			response.Ok(c)
		})
	}
}
