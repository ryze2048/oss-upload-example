package v1

import "oss-upload-example/api/v1/uploading"

type ApiGroup struct {
	UploadIngApiGroup uploading.UploadingApi
}

var ApiGroupApp = new(ApiGroup)
