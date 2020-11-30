package oss

import (
	"bytes"
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/astaxie/beego"
)

var ossClient *oss.Client

func init() {
	var (
		err error
	)
	endpoint := beego.AppConfig.String("oss::endpoint")
	// 阿里云主账号AccessKey拥有所有API的访问权限
	accessKeyId := beego.AppConfig.String("oss::access_key")
	accessKeySecret := beego.AppConfig.String("oss::access_secret")
	// 创建OSSClient实例。
	ossClient, err = oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		ossClient = nil
	}
}

//上传文件（路径方式）
func UploadByPath(bucketName string, localName string, saveName string) (string, error) {
	// 获取存储空间。
	if ossClient == nil {
		return "", errors.New("oss connect fail")
	}
	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		return "", errors.New("oss bucket not exist")
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(saveName, localName)
	if err != nil {
		return "", errors.New("upload fail")
	}
	returnUrl := beego.AppConfig.String("oss::bucket_domain") + saveName
	return returnUrl, nil
}

//上传图片 （Form上传）
func UploadByBuffer(bucketName string, buffer []byte, saveName string) (string, error) {
	// 获取存储空间。
	if ossClient == nil {
		return "", errors.New("oss connect fail")
	}
	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		return "", errors.New("oss bucket not exist")
	}
	// 上传文件
	options := []oss.Option{
		oss.ContentType("image/jpeg"),
	}
	err = bucket.PutObject(saveName, bytes.NewReader(buffer), options...)
	if err != nil {
		return "", errors.New("upload fail")
	}
	returnUrl := beego.AppConfig.String("oss::bucket_domain") + "/" + saveName
	return returnUrl, nil
}
