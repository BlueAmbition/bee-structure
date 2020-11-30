package base

import (
	"bee-structure/cache"
	"bee-structure/functions/array"
	"bee-structure/functions/jwt"
	"bee-structure/functions/redis"
	"bee-structure/functions/str"
	"bee-structure/models/user"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	//成功code
	SuccessCode = 200
	//创建成功代码
	CreatedCode = 201
	//失败code
	FailureCode = 500
	//需要授权
	UnAuthorized = 401
	//禁止访问，权限不够
	ForbiddenCode = 403
	//未找到
	NotFoundCode = 404
	//数据校验失败
	UnprocessableEntityCode = 422
	//有更新
	MovedPermanentlyCode = 301
	//无支付密码
	NoPayPasswordCode = 501
)

type BaseController struct {
	beego.Controller
}

/**
 * 通用返回Map对象
 *@param code 返回状态码。
 *@param msg 返回信息。
 *@param data 返回的数据对象。
 *@return map[string]interface{} 通用Map对象
 */
func (c *BaseController) GeneralReturn(code int, msg string, data interface{}) map[string]interface{} {
	mapData := map[string]interface{}{
		"code": code, "status": "failure", "msg": msg}
	//状态节点处理
	successArr := []string{strconv.Itoa(SuccessCode), strconv.Itoa(CreatedCode)}
	if array.InArray(strconv.Itoa(code), successArr) {
		mapData["status"] = "success"
	}
	//数据节点处理
	if data != nil {
		mapData["data"] = data
	}

	return mapData
}

/**
 * 统一返回Json格式
 *@param code 返回状态码。
 *@param msg 返回信息。
 *@return map[string]interface{} 通用Map对象
 */
func (c *BaseController) ResJson(code int, msg string, data interface{}) {
	c.Data["json"] = c.GeneralReturn(code, msg, data)
	c.ServeJSON()
}

/**
 * 从请求头获取Token令牌
 *@return string 获取的Token令牌
 */
func (c *BaseController) GetToken() string {
	return c.Ctx.Request.Header.Get("Authorization")
}

/**
 * 从Token中获取用户ID
 *@return int64 获取用户信息
 */
func (c *BaseController) GetTokenUserId() int64 {
	token := c.GetToken()
	if token == "" {
		return 0
	}
	flag, userId := jwt.CheckToken(token)
	if !flag {
		msg := fmt.Sprintf("用户ID：%v，请求地址：%v，请求头token：%v，客户端：%v", userId, c.Ctx.Input.URL(), token, c.UserAgent())
		logs.GetLogger("Login").Println(msg)
		c.Abort("401")
	}
	return userId
}

/**
 * 获取用户信息
 *@return custom.UserInfo 返回用户信息
 */
func (c *BaseController) GetUserInfo() user.UserInfo {
	userId := c.GetTokenUserId()
	user := cache.GetUserInfo(userId)
	if user.Id < 1 || user.Status != 1 {
		c.Abort("403")
	}
	return user
}

/**
 * 统一校验Json对象是否合法
 *@param objStruct 反序列化Json对象。
 *@return bool 验证结果
 */
func (c *BaseController) ValidateJSON(objStruct interface{}) bool {
	err := json.Unmarshal(c.Ctx.Input.RequestBody, objStruct)
	if err != nil {
		msg := c.ReturnMsg("data_error")
		c.ResJson(UnprocessableEntityCode, msg, nil)
		return false
	}
	return true
}

/**
 * 证书加密
 *@param originData 原始数据。
 *@return []byte 加密数据
 *@return error 错误信息
 */
func (c *BaseController) Encrypt(originData []byte) ([]byte, error) {
	publicKey, err := ioutil.ReadFile("public.pem")
	if err != nil {
		return nil, errors.New("failed load public key file")
	}
	encryptStr, err := str.RsaEncrypt(publicKey, originData)
	return encryptStr, err
}

/**
 * 从头信息或者url参数中获取当前语言
 *@return string 获取language_word的相应字段部分
 */
func (c *BaseController) GetLanguageKey() string {
	lang := c.Ctx.Request.Header.Get("Show-Language")
	if lang == "" {
		lang = c.GetString("lang", "")
		if lang == "" {
			lang = "zh-CN"
		}
	}
	//转换成language_word表中的相应的语言字段
	lang = strings.ToLower(lang)
	//过滤语言设置默认
	if lang != "zh-cn" && lang != "en-us" && lang != "zh-hk" {
		lang = "en-us"
	}
	lang = strings.Replace(lang, "-", "_", -1)
	return lang
}

/**
 * 根据key获取语言信息
 *@return string 返回信息
 */
func (c *BaseController) ReturnMsg(key string, args ...interface{}) string {
	lang := c.GetLanguageKey()
	word, error := cache.GetTipsWord(key, lang)
	if error != nil {
		return error.Error()
	}
	return fmt.Sprintf(word, args...)
}

/**
 * 允许注册国家
 *@return bool 返回信息
 *@return int 返回信息
 */
func (c *BaseController) AllowRegArea(areaCode string) (bool, int64) {
	//允许注册的国家
	allows := cache.GetCountryList()
	for _, v := range allows {
		if v.MobileCode == areaCode {
			return true, v.Id
		}
	}
	return false, 0
}

/**
 * 验证码校验
 *@param redisKey redis key。
 *@param validCode 校验验证码。
 *@return bool 返回信息
 */
func (c *BaseController) ValidateCode(redisKey, validCode string) bool {
	redisCode, err := redis.GetString(0, redisKey)
	if err != nil {
		return false
	}
	codeArr := strings.Split(redisCode, "|")
	if len(codeArr) != 2 {
		return false
	}
	code := codeArr[0]
	if !strings.EqualFold(code, validCode) {
		return false
	}
	return true
}

/**
 * 获取表行为字段文本
 *@param key behavior key。
 *@param behavior 行为值。
 *@param appendStr 拼接字符串
 *@return string 返回信息
 */
func (c *BaseController) BehaviorTxt(key string, behavior int, appendStr string) string {
	source := c.ReturnMsg("behavior_" + key)
	arr := strings.Split(source, "|")
	if behavior < 0 || behavior > len(arr)-1 {
		return c.ReturnMsg("unknown") // "未知"
	}
	if appendStr != "" {
		return arr[behavior] + appendStr
	}
	return arr[behavior]
}

//钱包行为
type WalletBehavior struct {
	Behavior    int    `json:"behavior"`
	BehaviorTxt string `json:"behavior_txt"`
}

/**
 * 用户钱包行文Json处理
 *@param key behavior key。
 *@param behavior 行为值。
 *@param appendStr 拼接字符串
 *@return string 返回信息
 */
func (c *BaseController) WalletBehaviorTxt(behavior int, appendStr string) string {
	var (
		list []WalletBehavior
	)
	source := c.ReturnMsg("behavior_user_wallet_record")
	err := json.Unmarshal([]byte(source), &list)
	if err != nil {
		return c.ReturnMsg("unknown") // "未知"
	}
	for _, v := range list {
		if v.Behavior == behavior {
			return v.BehaviorTxt
		}
	}
	return c.ReturnMsg("unknown")
}

/**
 * 图片等静态文件域名
 *@return string 返回信息
 */
func (c *BaseController) GetImgDomain() string {
	domain := beego.AppConfig.String("oss::domain")
	//OSS上层目录
	pathPre := beego.AppConfig.String("oss::path_pre")
	return domain + "/" + pathPre
}

/**
 * 获取表状态字段文本
 *@param key status key。
 *@param status 行为值。
 *@return string 返回信息
 */
func (c *BaseController) StatusTxt(key string, status int) string {
	source := c.ReturnMsg("status_" + key)
	arr := strings.Split(source, "|")
	if status < 0 || status > len(arr)-1 {
		return c.ReturnMsg("unknown")
	}
	return arr[status]
}

/**
 * 获取表类型字段文本
 *@param key type key。
 *@param status 状态。
 *@return string 返回信息
 */
func (c *BaseController) TypeTxt(key string, status int) string {
	source := c.ReturnMsg("type_" + key)
	arr := strings.Split(source, "|")
	if status < 0 || status > len(arr)-1 {
		return c.ReturnMsg("unknown") //"未知"
	}
	return arr[status]
}

/**
 * 从头信息或者url参数中获取当前语言Id
 *@return int LanguageId
 */
func (c *BaseController) GetLanguageId() int64 {
	lang := c.Ctx.Request.Header.Get("Show-Language")
	if lang == "" {
		lang = c.GetString("lang", "")
		if lang == "" {
			lang = "en-US"
		}
	}
	languageInfo := cache.GetLangByLocalKey(lang)
	if languageInfo.Id > 0 {
		return languageInfo.Id
	}
	return 2 //默认英语
}

//判断客户端
func (c *BaseController) UserAgent() string {
	userAgent := c.Ctx.Request.Header.Get("User-Agent")
	if userAgent != "" {
		userAgent = strings.ToLower(userAgent)
		if strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "ipad") {
			return "ios"
		}
		if strings.Contains(userAgent, "android") {
			return "android"
		}
		return userAgent
	}
	return ""
}
