// @APIVersion 1.0.0
// @Title ITI API
// @Description ITI APP相关API
// @Contact astaxie@gmail.com
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"bee-structure/controllers"
	"bee-structure/controllers/app"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	//app模块
	appMod := beego.NewNamespace("/app",
		beego.NSInclude(
			&app.AppController{},
		),
		//多国语言词组AppLang
		beego.NSRouter("/app-lang", &app.AppController{}, "get:AppLang"),
		//多国语言列表
		beego.NSRouter("/lang-list", &app.AppController{}, "get:LangList"),
		//更新检查
		beego.NSRouter("/upgrade", &app.AppController{}, "get:Upgrade"),
		//APP初始化配置信息
		beego.NSRouter("/init-configs", &app.AppController{}, "get:InitConfigs"),
	)

	ns := beego.NewNamespace("/v1",
		//index模块
		beego.NSNamespace("/index",
			//需要使用文档控制器
			beego.NSInclude(
				&controllers.IndexController{},
			),
			//登录
			beego.NSRouter("/login", &controllers.IndexController{}, "post:Login"),
			//注册
			beego.NSRouter("/register", &controllers.IndexController{}, "put:Register"),
			//退出
			beego.NSRouter("/logout", &controllers.IndexController{}, "post:Logout"),
			//允许注册的国家
			beego.NSRouter("/register-allow-country", &controllers.IndexController{}, "get:RegAllowCountry"),
			//找回密码
			beego.NSRouter("/find-password", &controllers.IndexController{}, "patch:FindPassword"),
			//图片上传
			//beego.NSRouter("/upload-img", &controllers.IndexController{}, "post:UploadImg"),
		),
	)
	beego.AddNamespace(appMod)
	beego.AddNamespace(ns)

	//跨域处理
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Version"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))
}
