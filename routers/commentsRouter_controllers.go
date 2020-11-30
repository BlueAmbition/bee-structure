package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["bee-structure/controllers:IndexController"] = append(beego.GlobalControllerRouter["bee-structure/controllers:IndexController"],
        beego.ControllerComments{
            Method: "FindPassword",
            Router: "/find-password",
            AllowHTTPMethods: []string{"patch"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bee-structure/controllers:IndexController"] = append(beego.GlobalControllerRouter["bee-structure/controllers:IndexController"],
        beego.ControllerComments{
            Method: "Index",
            Router: "/index",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bee-structure/controllers:IndexController"] = append(beego.GlobalControllerRouter["bee-structure/controllers:IndexController"],
        beego.ControllerComments{
            Method: "InitConfigs",
            Router: "/init-configs",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bee-structure/controllers:IndexController"] = append(beego.GlobalControllerRouter["bee-structure/controllers:IndexController"],
        beego.ControllerComments{
            Method: "Login",
            Router: "/login",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bee-structure/controllers:IndexController"] = append(beego.GlobalControllerRouter["bee-structure/controllers:IndexController"],
        beego.ControllerComments{
            Method: "Logout",
            Router: "/logout",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bee-structure/controllers:IndexController"] = append(beego.GlobalControllerRouter["bee-structure/controllers:IndexController"],
        beego.ControllerComments{
            Method: "Register",
            Router: "/register",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bee-structure/controllers:IndexController"] = append(beego.GlobalControllerRouter["bee-structure/controllers:IndexController"],
        beego.ControllerComments{
            Method: "RegAllowCountry",
            Router: "/register-allow-country",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bee-structure/controllers:IndexController"] = append(beego.GlobalControllerRouter["bee-structure/controllers:IndexController"],
        beego.ControllerComments{
            Method: "UploadImg",
            Router: "/upload-img",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
