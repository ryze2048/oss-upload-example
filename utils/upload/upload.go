package upload

import "oss-upload-example/global"

type OSS interface {
	GetSignUrl(key string) (string, error)
	GetPolicyToken() (*PolicyToken, error)
}

func NewOss() OSS {
	switch global.CONFIG.System.OssType {
	case "aliyun-oss":
		return &AliyunOSS{}
	default:
		return &AliyunOSS{}
	}
}
