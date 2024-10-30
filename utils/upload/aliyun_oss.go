package upload

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"hash"
	"io"
	"oss-upload-example/global"
	"time"
)

type AliyunOSS struct{}

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type PolicyToken struct {
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
	Callback    string `json:"callback"`
}

type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

func newClient() (client *oss.Client, err error) {
	return oss.New(global.CONFIG.AliyunOSS.Endpoint, global.CONFIG.AliyunOSS.AccessKeyId, global.CONFIG.AliyunOSS.AccessKeySecret)
}

func newBucket() (*oss.Bucket, error) {
	var err error
	var client *oss.Client
	if client, err = newClient(); err != nil {
		return nil, err
	}
	return client.Bucket(global.CONFIG.AliyunOSS.BucketName)
}

func (a *AliyunOSS) GetSignUrl(fileName string) (singUrl string, err error) {
	var bucket *oss.Bucket
	if bucket, err = newBucket(); err != nil {
		return "", err
	}

	return bucket.SignURL(fileName, oss.HTTPPut, 3600)
}

func (a *AliyunOSS) GetPolicyToken() (*PolicyToken, error) {
	var err error
	var now = time.Now().Unix()
	var expireEnd = now + global.CONFIG.AliyunOSS.ExpireTime
	var tokenExpire = getGmtIso8601(expireEnd)

	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, global.CONFIG.AliyunOSS.BasePath)
	config.Conditions = append(config.Conditions, condition)

	var result = make([]byte, 0)
	if result, err = json.Marshal(config); err != nil {
		return nil, err
	}
	deByte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(global.CONFIG.AliyunOSS.AccessKeySecret))
	if _, err = io.WriteString(h, deByte); err != nil {
		return nil, err
	}
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var callbackParam CallbackParam
	callbackParam.CallbackUrl = global.CONFIG.AliyunOSS.CallbackUrl
	callbackParam.CallbackBody = "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"

	var callbackStr = make([]byte, 0)
	if callbackStr, err = json.Marshal(callbackParam); err != nil {
		return nil, err
	}
	callbackBase64 := base64.StdEncoding.EncodeToString(callbackStr)

	var policyToken PolicyToken
	policyToken.AccessKeyId = global.CONFIG.AliyunOSS.AccessKeyId
	policyToken.Host = global.CONFIG.AliyunOSS.Endpoint
	policyToken.Expire = expireEnd
	policyToken.Signature = signedStr
	policyToken.Directory = global.CONFIG.AliyunOSS.BasePath
	policyToken.Policy = deByte
	policyToken.Callback = callbackBase64

	return &policyToken, nil
}

func getGmtIso8601(expireEnd int64) string {
	var tokenExpire = time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")
	return tokenExpire
}
