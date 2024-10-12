package config

type Server struct {
	AliyunOSS AliyunOSS `mapstructure:"aliyun-oss" json:"aliyun-oss" yaml:"aliyun-oss"`
	System    System    `mapstructure:"system" json:"system" yaml:"system"`
}
