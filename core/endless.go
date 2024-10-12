package core

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"time"
)

func initEndless(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 360 * time.Second
	s.WriteTimeout = 360 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}
