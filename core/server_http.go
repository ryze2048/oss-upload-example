package core

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"oss-upload-example/global"
	"oss-upload-example/initialize"
	"oss-upload-example/utils"
	"time"
)

type server interface {
	ListenAndServe() error
}

func RunHttpServer() {
	utils.NormalDaemon(entry)
}

func entry(ctx context.Context) {

	Router := initialize.Routers() // 路由初始化

	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)

	log.Info("service run success on ", address)

	s := initEndless(address, Router)
	go func() {
		<-ctx.Done()
	}()
	if err := s.ListenAndServe(); err != nil {
		log.Error("listen service err --> ", err)
	}
	log.Info("service stop ", time.Now().Format("2006-01-02 15:04:05"))
}
