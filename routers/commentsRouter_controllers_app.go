package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["bee-structure/controllers/app:AppController"] = append(beego.GlobalControllerRouter["bee-structure/controllers/app:AppController"],
        beego.ControllerComments{
            Method: "AppLang",
            Router: "/app-lang",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bee-structure/controllers/app:AppController"] = append(beego.GlobalControllerRouter["bee-structure/controllers/app:AppController"],
        beego.ControllerComments{
            Method: "InitConfigs",
            Router: "/init-configs",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bee-structure/controllers/app:AppController"] = append(beego.GlobalControllerRouter["bee-structure/controllers/app:AppController"],
        beego.ControllerComments{
            Method: "LangList",
            Router: "/lang-list",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bee-structure/controllers/app:AppController"] = append(beego.GlobalControllerRouter["bee-structure/controllers/app:AppController"],
        beego.ControllerComments{
            Method: "Upgrade",
            Router: "/upgrade",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
