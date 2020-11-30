package base

import (
	"bee-structure/cache"
	"bee-structure/functions/redis"
	"bee-structure/functions/req"
	"bee-structure/functions/str"
	"bee-structure/models/blockchain"
	"bee-structure/models/user"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserBaseController struct {
	BaseController
}

//预处理函数
func (c *UserBaseController) Prepare() {
	token := c.GetToken()
	userId := c.GetTokenUserId()
	tokenStr := strings.Replace(token, "Bearer ", "", 1)
	redisKey := fmt.Sprint("login-token:", userId)
	redisToken, err := redis.GetString(0, redisKey)
	if err != nil || redisToken == "" || redisToken != tokenStr {
		msg := fmt.Sprintf("用户ID：%v，请求地址：%v，请求头token：%v，redis token：%v，客户端：%v", userId, c.Ctx.Input.URL(), token, redisToken, c.UserAgent())
		logs.GetLogger("Login").Println(msg)
		c.Abort("401")
		return
	}
}

//验证交易密码
func (c *UserBaseController) ValidatePayPassWord(payPassWord string) bool {
	userId := c.GetUserInfo().Id
	userInfo := user.GetUserById(userId)
	if userInfo.PayPassword == "" {
		msg := c.ReturnMsg("no_pay_password")
		c.ResJson(NoPayPasswordCode, msg, nil)
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(userInfo.PayPassword), []byte(payPassWord))
	if err != nil {
		msg := c.ReturnMsg("user_pay_password_error")
		c.ResJson(UnprocessableEntityCode, msg, nil)
		return false
	}
	return true
}

/**
 * 根据币种获取区块链域名
 *@param coinKey 币种Key。
 *@return chain 所属链
 *@return string 返回配置请求对应域名
 */
func (c *BaseController) BlockChainDomain(coinKey string, chain int64) string {
	coinKey = strings.ToLower(coinKey)
	domain := ""
	switch coinKey {
	case "eth":
		domain = beego.AppConfig.String("ethereum::domain")
		break
	case "usdt":
		domain = beego.AppConfig.String("omni_usdt::domain")
		if chain == 1 {
			domain = beego.AppConfig.String("erc20_usdt::domain")
		} else if chain == 2 {
			domain = beego.AppConfig.String("tt_usdt::domain")
		}
		break
	case "btb":
		domain = beego.AppConfig.String("btb::domain")
		break

	}
	return domain
}

/**
 * 获取币种加密公钥Key
 *@param coinKey 币种Key。
 *@return chain 所属链
 *@return string 返回配置请求对应域名
 */
func (c *BaseController) BlockChainPublicKey(coinKey string) string {
	coinKey = strings.ToLower(coinKey)
	switch coinKey {
	case "btc":
		return "btc_rpc_pk"
	case "eth":
		return "eth_rpc_pk"
	case "usdt":
		return "omnicore_rpc_pk"
	case "eos":
		return "eos_rpc_pk"
	case "btb":
		return "btb_rpc_pk"
	}
	return ""
}

/**
 * 请求区块链接口
 *@param platformKey 币种Key。
 *@return chain 所属链
 *@return bool 请求是否返回成功
 *@return custom.BlockChainRes 区块链通用返回信息
 */
func (c *BaseController) ReqBlockChain(platformKey string, chain int64, address string, mapParams map[string]interface{}) (bool, blockchain.BlockChainRes) {
	domain := c.BlockChainDomain(platformKey, chain)
	var resObj blockchain.BlockChainRes
	if domain == "" {
		msg := fmt.Sprintf("调用域名为空，币种key： %v", platformKey)
		logs.Error("[BlockChain]", msg)
		return false, resObj
	}
	bcKey := c.BlockChainPublicKey(platformKey)
	if bcKey == "" {
		msg := fmt.Sprintf("获取公钥key失败，币种key： %v", platformKey)
		logs.Error("[BlockChain]", msg)
		return false, resObj
	}
	paramJson, _ := json.Marshal(mapParams)
	url := domain + address
	mapData := make(map[string][]string)
	config := cache.GetBlockChainConfig(bcKey)
	publicKey := []byte(config.Value)
	encryptStr, err := str.RsaEncrypt(publicKey, paramJson)
	if err != nil {
		msg := fmt.Sprintf("参数加密失败:%v，原始数据：%v，请求地址:%v，币种key： %v", err.Error(), paramJson, url, platformKey)
		logs.Error("[BlockChain]", msg)
		return false, resObj
	}
	mapData["data"] = []string{string(encryptStr)}
	res, err := req.PostForm(url, mapData)
	if err != nil {
		msg := fmt.Sprintf("请求失败:%v，原始数据：%v，请求地址:%v，币种key： %v", err.Error(), paramJson, url, platformKey)
		logs.Error("[BlockChain]", msg)
		return false, resObj
	}
	err = json.Unmarshal(res, &resObj)
	if err != nil || resObj.Code != SuccessCode {
		var msg string
		if err != nil {
			msg = fmt.Sprintf("返回数据错误:%v，原始数据：%v，请求地址:%v，币种key： %v，返回数据：%v", err.Error(), paramJson, url, platformKey, string(res))
		} else {
			msg = fmt.Sprintf("原始数据：%v，请求地址:%v，币种key： %v，返回数据：%v", paramJson, url, platformKey, string(res))
		}
		logs.Error("[BlockChain]", msg)
		return false, resObj
	}
	return true, resObj
}
