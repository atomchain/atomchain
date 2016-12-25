package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["chain/net/rpc/controllers:ChainController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:ChainController"],
        beego.ControllerComments{
            Method: "CreateAccount",
            Router: `/createaccount`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chain/net/rpc/controllers:ChainController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:ChainController"],
        beego.ControllerComments{
            Method: "RemoveAccount",
            Router: `/removeaccount`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chain/net/rpc/controllers:ChainController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:ChainController"],
        beego.ControllerComments{
            Method: "SendTransaction",
            Router: `/sendtransaction`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chain/net/rpc/controllers:ChainController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:ChainController"],
        beego.ControllerComments{
            Method: "SendVote",
            Router: `/sendvote`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chain/net/rpc/controllers:ObjectController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chain/net/rpc/controllers:ObjectController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chain/net/rpc/controllers:ObjectController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chain/net/rpc/controllers:ObjectController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chain/net/rpc/controllers:ObjectController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chain/net/rpc/controllers:UserController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:UserController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chain/net/rpc/controllers:UserController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chain/net/rpc/controllers:UserController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:UserController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chain/net/rpc/controllers:UserController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:UserController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chain/net/rpc/controllers:UserController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:UserController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chain/net/rpc/controllers:UserController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:UserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["chain/net/rpc/controllers:UserController"] = append(beego.GlobalControllerRouter["chain/net/rpc/controllers:UserController"],
        beego.ControllerComments{
            Method: "Logout",
            Router: `/logout`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
