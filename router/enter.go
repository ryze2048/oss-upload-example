package router

import "oss-upload-example/router/uploading"

type RouterGroup struct {
	Uploading uploading.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
