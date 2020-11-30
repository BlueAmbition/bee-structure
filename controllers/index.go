package controllers

import (
	"bee-structure/cache"
	"bee-structure/controllers/base"
	"bee-structure/functions/jwt"
	"bee-structure/functions/redis"
	"bee-structure/functions/str"
	"bee-structure/functions/valid"
	"bee-structure/models/user"
	"encoding/json"
	"github.com/astaxie/beego"
	"golang.org/x/crypto/bcrypt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type IndexController struct {
	base.BaseController
}

// @Title 通用返回说明
// @Description 通用返回说明
// @Success 200 {"code":200,"status":"success","msg":"注册成功"}
// @Failure 422 {"code":422,"status":"failure","msg":"参数有误"}
// @Failure 500 {"code":500,"status":"failure","msg":"创建失败"}
// @Failure 401 {"code":401,"status":"failure","msg":"登陆失效"}
// @Failure 403 {"code":403,"status":"failure","msg":"用户信息不存在或被禁止登陆"}
// @Success 301 {"code":301,"status":"success","msg":"需要强制升级","data":{"id":"int64 主键","app_name":"string app名","version":"string 版本号","intro":"string 简介","force":"int 1:强制升级","android_url":"string 安卓包下载地址","ios_url":"string ios包下载地址","need_upgrade":"bool 是否需要升级","title": "string 升级标题","content": "string 升级内容","created_at":"2019-09-24T14:43:38+08:00","updated_at":"2020-01-08T15:07:02+08:00"}}
// @router /index [get]
func (c *IndexController) Index() {

}

// @Title 用户登陆
// @Description 用户登陆，手机或邮箱
// @Param user_name  formData 	string	 true	"用户名email或mobile"
// @Param password  formData 	string	 true	"密码"
// @Success 200 {"code":200,"data":{"auth_token":"string APP授权Token","im_token":"string IM授权Token"},"msg":"登录成功","status":"success"}
// @router /login [post]
func (c *IndexController) Login() {
	var (
		msg       string
		err       error
		objStruct struct {
			UserName string `json:"user_name"`
			Password string `json:"password"`
		}
	)
	if !c.ValidateJSON(&objStruct) {
		return
	}
	user := user.GetUserByUserName(objStruct.UserName)
	if user.Id <= 0 {
		msg = c.ReturnMsg("failure_index_Login")
		c.ResJson(base.UnprocessableEntityCode, msg, nil)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(objStruct.Password))
	if err != nil {
		msg = c.ReturnMsg("failure_index_Login")
		c.ResJson(base.UnprocessableEntityCode, msg, nil)
		return
	}
	var expireSeconds int64 = 2592000
	token := jwt.MakeToken(user.Id, expireSeconds)
	if token == "" {
		msg = c.ReturnMsg("failure_index_Login")
		c.ResJson(base.UnprocessableEntityCode, msg, nil)
		return
	}
	//token过期时间设置成30天
	flag := redis.SetString(0, "login-token:"+strconv.FormatInt(user.Id, 10), token, int(expireSeconds))
	if !flag {
		msg = c.ReturnMsg("failure_index_Login")
		c.ResJson(base.FailureCode, msg, nil)
		return
	}

	//tag := user.Nickname
	//imToken := c.IMLogin(user.Id, user.HeadImg, tag)
	//redis.SetString(0, "im-token:"+strconv.FormatInt(user.Id, 10), imToken, int(expireSeconds))
	cache.SetUserInfo(user.Id)
	data := map[string]string{"auth_token": token}
	msg = c.ReturnMsg("success_login")
	c.ResJson(base.SuccessCode, msg, data)
}

// @Title 用户注册
// @Description 手机号或邮箱注册用户
// @Param type  formData 	string	 true	"注册类型,email或mobile"
// @Param mobile  formData 	string	 true	"手机号，邮箱注册传空"
// @Param email  formData 	string	 true	"邮箱，手机注册传空"
// @Param nick_name  formData 	string	 true	"昵称"
// @Param login_password  formData 	string	 true	"登录密码"
// @Param invite_code  formData 	string	 true	"邀请码"
// @Param valid_code  formData 	string	 true	"验证码"
// @Param area_code  formData 	string	 true	"地区编号"
// @Success 200 {"code":200,"status":"success","msg":"注册成功"}
// @router /register [put]
func (c *IndexController) Register() {
	var (
		msg       string
		flag      bool
		redisKey  string
		countryId int64
		objStruct struct {
			RegType       string `json:"type"`
			Email         string `json:"email"`
			NickName      string `json:"nick_name"`
			Mobile        string `json:"mobile"`
			LoginPassword string `json:"login_password"`
			//PayPassword   string `json:"pay_password"`
			ValidCode  string `json:"valid_code"`
			InviteCode string `json:"invite_code"`
			AreaCode   string `json:"area_code"`
			//CountryId     int    `json:"country_id"`
		}
		flagChan chan bool
	)
	if !c.ValidateJSON(&objStruct) {
		return
	}
	if str.StringLen(objStruct.NickName) < 2 || str.StringLen(objStruct.NickName) > 8 {
		msg = c.ReturnMsg("register_nickname_error")
		c.ResJson(base.UnprocessableEntityCode, msg, nil)
		return
	}

	if len(objStruct.LoginPassword) < 8 {
		msg = c.ReturnMsg("register_password_error", 8)
		c.ResJson(base.UnprocessableEntityCode, msg, nil)
		return
	}

	//if len(objStruct.PayPassword) < 8 {
	//	msg = c.ReturnMsg("register_pay_password_error", 8)
	//	c.ResJson(base.UnprocessableEntityCode, msg, nil)
	//	return
	//}

	//验证码处理
	objStruct.ValidCode = strings.Trim(objStruct.ValidCode, " ")
	if objStruct.ValidCode == "" {
		msg = c.ReturnMsg("register_code_error")
		c.ResJson(base.UnprocessableEntityCode, msg, nil)
		return
	}

	//注册类型
	if objStruct.RegType == "email" {
		//邮箱验证方式，邮箱是否存在
		objStruct.Email = strings.Trim(objStruct.Email, " ")
		flag = valid.IsEmail(objStruct.Email)
		if !flag {
			msg = c.ReturnMsg("email_error")
			c.ResJson(base.UnprocessableEntityCode, msg, nil)
			return
		}
		flag = user.IsEmailExist(objStruct.Email)
		if flag {
			msg = c.ReturnMsg("register_email_error")
			c.ResJson(base.UnprocessableEntityCode, msg, nil)
			return
		}
		redisKey = "reg_email_" + objStruct.Email
		flag = c.ValidateCode(redisKey, objStruct.ValidCode)
		if !flag {
			msg, _ = cache.GetTipsWord("code_error", c.GetLanguageKey())
			c.ResJson(base.UnprocessableEntityCode, msg, nil)
			return
		}

	} else if objStruct.RegType == "mobile" {
		//允许注册的国家
		flag, countryId = c.AllowRegArea(objStruct.AreaCode)
		if !flag {
			msg = c.ReturnMsg("reg_allow_area")
			c.ResJson(base.UnprocessableEntityCode, msg, nil)
			return
		}
		//手机验证方式，手机号是否存在
		objStruct.Mobile = strings.Trim(objStruct.Mobile, " ")
		flag = valid.IsMobile(objStruct.Mobile, objStruct.AreaCode)
		if !flag {
			msg = c.ReturnMsg("mobile_error")
			c.ResJson(base.UnprocessableEntityCode, msg, nil)
			return
		}
		flag = user.IsMobileExist(objStruct.Mobile)
		if flag {
			msg = c.ReturnMsg("register_mobile_error")
			c.ResJson(base.UnprocessableEntityCode, msg, nil)
			return
		}
		unionMobile := objStruct.AreaCode + objStruct.Mobile
		redisKey = "reg_mobile_" + unionMobile
		flag = c.ValidateCode(redisKey, objStruct.ValidCode)
		if !flag {
			msg, _ = cache.GetTipsWord("code_error", c.GetLanguageKey())
			c.ResJson(base.UnprocessableEntityCode, msg, nil)
			return
		}
	} else {
		msg = c.ReturnMsg("register_type_error")
		c.ResJson(base.UnprocessableEntityCode, msg, nil)
		return
	}
	if objStruct.InviteCode == "" {
		msg = c.ReturnMsg("register_invite_code_error")
		c.ResJson(base.UnprocessableEntityCode, msg, nil)
		return
	}
	//邀请码处理
	objStruct.InviteCode = strings.Trim(objStruct.InviteCode, "")
	user1 := user.User{Mobile: objStruct.Mobile, Email: objStruct.Email, CountryId: countryId, Nickname: objStruct.NickName,
		MobileCode: objStruct.AreaCode, Password: objStruct.LoginPassword,
		ParentId: 0, ParentTree: "",
	}
	parentUser := user.GetUserByInviteCode(objStruct.InviteCode)
	if parentUser.Id > 0 {
		user1.ParentId = parentUser.Id
		user1.ParentTree = parentUser.ParentTree + "," + strconv.FormatInt(parentUser.Id, 10) + ","
		user1.ParentTree = strings.Trim(user1.ParentTree, ",")
		user1.Level = parentUser.Level + 1
	} else {
		msg = c.ReturnMsg("invite_code_error")
		c.ResJson(base.UnprocessableEntityCode, msg, nil)
		return
	}
	id, _ := user.Register(objStruct.RegType, user1)
	if id > 0 {
		redis.DelKey(0, redisKey)
		flagChan = make(chan bool, 2)
		//生成邀请码
		go initInviteCode(id, flagChan)
		msg = c.ReturnMsg("register_success")
		c.ResJson(base.SuccessCode, msg, nil)
		return
	}
	msg = c.ReturnMsg("request_fail")
	c.ResJson(base.FailureCode, msg, nil)
}

/**
 * 生成注册用户邀请码
 *@param userId 用户ID。
 *@param flagChain channel。
 */
func initInviteCode(userId int64, flagChain chan bool) {
	user.GetInviteCode(userId)
	flagChain <- true
}

// @Title 用户退出
// @Description 退出成功
// @Success 200 {"code":200,"status":"success","msg":"退出成功"}
// @router /logout [post]
func (c *IndexController) Logout() {
	userInfo := c.GetUserInfo()
	idStr := strconv.FormatInt(userInfo.Id, 10)
	redis.DelKey(0, "login-token:"+idStr)
	redis.DelKey(0, "im-token:"+idStr)
	redis.DelKey(0, "user_info:"+idStr)
	msg := c.ReturnMsg("logout_success")
	c.ResJson(base.SuccessCode, msg, nil)
}

// @Title 允许注册的国家列表
// @Description 允许注册的国家列表
// @Success 200 {"code":200,"data":[{"id":"int 国家ID","country":"string 国家","country_en":"string 国家英文","mobile_code":"string 国家地区码"}],"msg":"Successful data acquisition.","status":"success"}
// @router /register-allow-country [get]
func (c *IndexController) RegAllowCountry() {
	allows := cache.GetCountryList()
	msg := c.ReturnMsg("success_get_data")
	c.ResJson(base.SuccessCode, msg, allows)
}

// @Title 找回密码
// @Description 登录页找回密码
// @Param type  formData 	string	 true	"找回类型,email或mobile"
// @Param mobile  formData 	string	 true	"手机号，邮箱找回传空"
// @Param email  formData 	string	 true	"邮箱，手机找回传空"
// @Param password  formData 	string	 true	"登录密码"
// @Param valid_code  formData 	string	 true	"验证码"
// @Param area_code  formData 	string	 true	"地区编号"
// @Success 200 {"code":200,"status":"success","msg":"找回成功"}
// @router /find-password [patch]
func (c *IndexController) FindPassword() {
	var (
		msg       string
		flag      bool
		redisKey  string
		objStruct struct {
			FindType  string `json:"type"`
			Email     string `json:"email"`
			Mobile    string `json:"mobile"`
			AreaCode  string `json:"area_code"`
			ValidCode string `json:"valid_code"`
			Password  string `json:"password"`
		}
	)

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &objStruct)
	if err != nil {
		msg = c.ReturnMsg("data_error")
		c.ResJson(base.UnprocessableEntityCode, msg, nil)
		return
	}

	if objStruct.FindType == "email" {
		//邮箱验证方式，邮箱是否存在
		objStruct.Email = strings.Trim(objStruct.Email, " ")
		flag = valid.IsEmail(objStruct.Email)
		if !flag {
			msg = c.ReturnMsg("email_error")
			c.ResJson(base.UnprocessableEntityCode, msg, nil)
			return
		}
		flag = user.IsEmailExist(objStruct.Email)
		if !flag {
			msg = c.ReturnMsg("email_not_registered")
			c.ResJson(base.UnprocessableEntityCode, msg, nil)
			return
		}
		redisKey = "find_email_" + objStruct.Email
		if !c.ValidateCode(redisKey, objStruct.ValidCode) {
			msg, _ = cache.GetTipsWord("code_error", c.GetLanguageKey())
			c.ResJson(base.UnprocessableEntityCode, msg, nil)
			return
		}
	} else if objStruct.FindType == "mobile" {
		//手机验证方式，手机号是否存在
		objStruct.Mobile = strings.Trim(objStruct.Mobile, " ")
		flag = valid.IsMobile(objStruct.Mobile, objStruct.AreaCode)
		if !flag {
			msg = c.ReturnMsg("mobile_error")
			c.ResJson(base.UnprocessableEntityCode, msg, nil)
			return
		}

		flag = user.IsMobileExist(objStruct.Mobile)
		if !flag {
			msg = c.ReturnMsg("mobile_not_registered")
			c.ResJson(base.UnprocessableEntityCode, msg, nil)
			return
		}
		unionMobile := objStruct.AreaCode + objStruct.Mobile
		redisKey = "find_mobile_" + unionMobile
		if !c.ValidateCode(redisKey, objStruct.ValidCode) {
			msg, _ = cache.GetTipsWord("code_error", c.GetLanguageKey())
			c.ResJson(base.UnprocessableEntityCode, msg, nil)
			return
		}
	} else {
		msg = c.ReturnMsg("data_error")
		c.ResJson(base.UnprocessableEntityCode, msg, nil)
		return
	}
	if len(objStruct.Password) < 8 {
		msg = c.ReturnMsg("password_len_error", 8)
		c.ResJson(base.UnprocessableEntityCode, msg, nil)
		return
	}
	user1 := user.User{Mobile: objStruct.Mobile, Email: objStruct.Email,
		MobileCode: objStruct.AreaCode, Password: objStruct.Password}

	_, f := user.FindPassword(objStruct.FindType, user1)
	if f {
		redis.DelKey(0, redisKey)
		msg = c.ReturnMsg("request_success")
		c.ResJson(base.SuccessCode, msg, nil)
		return
	}
	msg = c.ReturnMsg("request_fail")
	c.ResJson(base.FailureCode, msg, nil)
	return
}

// @Title 图片上传
// @Description 图片上传
// @Param file  formData 	file	 true	"图片文件"
// @Success 200 {"code":200,"status":"success","msg":"上传成功"}
// @router /upload-img [post]
func (c *IndexController) UploadImg() {
	var (
		msg string
		err error
	)
	//必须登录才有上传权限，后续创建目录使用ID
	c.GetUserInfo()
	file, info, err := c.GetFile("file") //返回文件，文件信息头，错误信息
	if err != nil {
		msg = c.ReturnMsg("data_error")
		c.ResJson(base.FailureCode, msg, nil)
		return
	}
	defer file.Close()        //关闭上传的文件，否则出现临时文件不清除的情况
	fileName := info.Filename //将文件信息头的信息赋值给filename变量
	t := time.Now()
	paths := "static/upload/" + strconv.Itoa(t.Year()) + "/" + t.Month().String() + "/" + strconv.Itoa(t.Day())
	if _, err = os.Stat(paths); err != nil {
		err = os.MkdirAll(paths, 0711)
		if err != nil {
			msg = c.ReturnMsg("make_dir_error")
			c.ResJson(base.FailureCode, msg, nil)
			return
		}
	}
	newFileName := strconv.FormatInt(t.UnixNano()/1e6, 10) + path.Ext(fileName)
	paths = path.Join(paths, "/", strconv.Itoa(t.Day())+"Day"+newFileName)
	err = c.SaveToFile("file", paths) //保存文件的路径。保存在paths中   （文件名）
	if err != nil {
		msg = c.ReturnMsg("upload_file_fail")
		c.ResJson(base.FailureCode, msg, nil)
		return
	}
	msg = c.ReturnMsg("upload_file_success")
	c.ResJson(base.SuccessCode, msg, map[string]interface{}{"filename": "/" + paths})
}

// @Title APP初始化配置信息
// @Description APP初始化配置信息
// @Success 200 {"code":200,"data":{"im_app_key":"string IM AppKey","im_secret":"string IM Secret","news_link":"分享地址","img_domain":"string 图片域名","agreement_dir":"string 用户协议在线目录","words_refresh_ts":"多语言刷新时间戳"},"msg":"Successful data acquisition","status":"success"}
// @router /init-configs [get]
func (c *IndexController) InitConfigs() {
	mapReturn := make(map[string]interface{})
	domain := c.GetImgDomain()
	imAppKey := beego.AppConfig.String("im::app_key")
	imSecret := beego.AppConfig.String("im::secret")
	newsLink := beego.AppConfig.String("url::news_link")
	agreementDir := beego.AppConfig.String("url::agreement")
	mapReturn["news_link"] = newsLink
	mapReturn["img_domain"] = domain
	mapReturn["agreement_dir"] = agreementDir
	mapReturn["im_app_key"] = imAppKey
	mapReturn["im_secret"] = imSecret
	mapReturn["words_refresh_ts"] = cache.GetWordsRefreshTs()
	msg := c.ReturnMsg("success_get_data")
	c.ResJson(base.SuccessCode, msg, mapReturn)
}
