package global

import (
	"github.com/spf13/viper"
	"oss-upload-example/config"
)

var (
	VIPER  *viper.Viper
	CONFIG config.Server
)
