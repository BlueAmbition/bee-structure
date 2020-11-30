package sms

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"strings"
)

//发送短信
func SendYunSms(mobile string, sign string, templateId string, content []string) bool {
	api := beego.AppConfig.String("yun_sms::api")
	accessKey := beego.AppConfig.String("yun_sms::access_key")
	secret := beego.AppConfig.String("yun_sms::secret")
	//sign := beego.AppConfig.String("yun_sms::sign")
	//templateId := beego.AppConfig.String("yun_sms::template_id")
	req := httplib.Post(api)
	//post请求
	req.Param("accesskey", accessKey)
	req.Param("secret", secret)
	req.Param("sign", sign)
	req.Param("mobile", mobile)
	req.Param("templateId", templateId)
	if len(content) != 0 {
		replaceStr := ""
		for _, v := range content {
			replaceStr += v + "##"
		}
		replaceStr = strings.Trim(replaceStr, "##")
		req.Param("content", replaceStr)
	}
	response, err := req.String()
	if err != nil {
		//msg := fmt.Sprintf("手机号:%v；短信动态参数:%v；错误信息:%v", mobile, content, err.Error())
		//log.WriteLog("sms.log", msg, "Error")
		return false
	}
	var objStruct struct {
		Code    string `json:"code"`
		Msg     string `json:"msg"`
		BatchId string `json:"batchId"`
	}
	err2 := json.Unmarshal([]byte(response), &objStruct)
	if err2 != nil {
		//msg := fmt.Sprintf("手机号:%v；短信动态参数:%v；错误信息:%v", mobile, content, err2.Error())
		//log.WriteLog("sms.log", msg, "Error")
		return false
	}
	if objStruct.Code != "0" {
		//msg := fmt.Sprintf("手机号:%v；短信动态参数:%v；返回状态码：%v；返回信息：%v；返回批次：%v", mobile, content, objStruct.Code, objStruct.Msg, objStruct.BatchId)
		//log.WriteLog("sms.log", msg, "Error")
		return false
	}
	//msg := fmt.Sprintf("手机号:%v；短信动态参数:%v；返回状态码：%v；返回信息：%v；返回批次：%v", mobile, content, objStruct.Code, objStruct.Msg, objStruct.BatchId)
	//log.WriteLog("sms.log", msg, "Error")
	return true
}

//发送验证码短信
func SendSmsCode(mobile string, content []string) bool {

	sign := beego.AppConfig.String("yun_sms::sign")
	templateId := beego.AppConfig.String("yun_sms::template_id")
	flag := SendYunSms(mobile, sign, templateId, content)
	return flag
}
