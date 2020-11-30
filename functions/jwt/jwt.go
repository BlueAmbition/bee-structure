package jwt

import (
	_ "github.com/astaxie/beego"
	_ "github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

//生成Token
func MakeToken(id int64, expireSeconds int64) string {
	signKey := []byte("aaaAAA111~!#$%^&*")
	// Create the Claims
	claims := &jwt.MapClaims{
		"user_id": id,
		"nbf":     int64(time.Now().Unix()),
		"exp":     int64(time.Now().Unix() + expireSeconds),
		"iss":     "ITI",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString(signKey)
	if err == nil {
		return str
	}
	return ""
}

// 校验token是否有效
func CheckToken(tokenStr string) (bool, int64) {
	if tokenStr == "" {
		return false, 0
	}
	tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("aaaAAA111~!#$%^&*"), nil
	})
	if err!=nil{
		return false, 0
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		return true, int64(claims["user_id"].(float64))
	}
	return false, 0
}
