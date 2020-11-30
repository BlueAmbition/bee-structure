package main

import (
	"bee-structure/cache"
	error2 "bee-structure/controllers/error"
	"bee-structure/filters"
	_ "bee-structure/routers"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//mysql连接初始化
func mysqlInit() {
	host := beego.AppConfig.String("mysql::host")
	port := beego.AppConfig.String("mysql::port")
	user := beego.AppConfig.String("mysql::user")
	password := beego.AppConfig.String("mysql::password")
	db := beego.AppConfig.String("mysql::db")
	//注册mysql Driver
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.DefaultTimeLoc = time.Local
	//用户名:密码@tcp(url地址)/数据库
	conn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + db + "?charset=utf8mb4&loc=Local"
	//注册数据库连接
	orm.RegisterDataBase("default", "mysql", conn)
	//orm.Debug = true
}

//日志初始化
func logInit() {
	logs.SetLogFuncCall(true)
	//异步输出
	logs.Async()
	fileName := strings.Replace(beego.AppPath, "\\", "/", -1) + "/logs/project.log"
	config := fmt.Sprintf(`{"filename":"%v","level":%v}`, fileName, logs.LevelWarning)
	logs.SetLogger(logs.AdapterFile, config)
}

//beego配置初始化
func beegoInit() {
	//不需要渲染模板
	beego.BConfig.WebConfig.AutoRender = false
	//错误友好提示
	beego.ErrorController(&error2.ErrorController{})
	//自定义处理异常
	beego.BConfig.RecoverFunc = func(ctx *context.Context) {
		err := recover()
		if err != nil {
			//写入日志
			var (
				stack string
				msg   string
			)

			reqData := fmt.Sprintf("请求路径：%v；\n错误信息：%v ；\n请求头：%v；\n请求参数：%v，\n请求body：%v", ctx.Input.URL(), err, ctx.Request.Header, ctx.Request.Form, string(ctx.Input.RequestBody))
			logs.Critical(reqData)

			for i := 1; ; i++ {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				logs.Critical(fmt.Sprintf("%s:%d", file, line))
				stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
			}
			//错误友好返回
			lang := ctx.Input.Header("Show-Language")
			if lang == "" {
				lang = "en-US"
			}
			lang = strings.ToLower(strings.Replace(lang, "-", "_", -1))
			switch err {
			case "401":
				msg, _ = cache.GetTipsWord("error_401", lang)
				break
			case "403":
				msg, _ = cache.GetTipsWord("error_403", lang)
				break
			default:
				msg, _ = cache.GetTipsWord("error_500", lang)
				break
			}
			code := "500"
			if err == "401" {
				code = "401"
				//ctx.ResponseWriter.WriteHeader(401)
			}
			if err == "403" {
				code = "403"
			}
			ctx.Output.Header("content-type", "application/json")
			ctx.WriteString(`{"code":` + code + `,"msg":"` + msg + `","status":false}`)
		}
	}
}

//过滤器初始化
func filterInit() {
	filters.AppUpgrade()
}

func init() {
	logInit()
	mysqlInit()
	filterInit()
	beegoInit()
}

func main() {
	//文档生成配置
	// if beego.BConfig.RunMode == "dev" {
	// 	beego.BConfig.WebConfig.DirectoryIndex = true
	// 	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	// }

	beego.Run()
}
