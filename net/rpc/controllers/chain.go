package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
)

type TX struct {
	From   string
	To     string
	Amount int64
	Cast   string
}

type VT struct {
	From   string
	To     string
	Amount int64
	Cast   string
}

type Account struct {
	Method   string
	User     string
	Length   int
	Password string
}

/*
example

curl -X GET http://localhost:8080/v1/chain/createaccount?password=100&user=bigarm

curl -H "Content-Type:application/json" -X POST --data '{"From" :"0x10","To":"0x20","Amount":1024,"Cast"  :"VOTE"}' "http://localhost:8080/v1/chain/sendvote"

*/

// ChainController Operations about Chain
type ChainController struct {
	beego.Controller
}

// @router /createaccount [get]
func (u *ChainController) CreateAccount() {
	fmt.Println("ChainController.CreateAccount")
	password := u.GetString("password")
	user := u.GetString("user")
	fmt.Printf("password=%s, user=%s\n", password, user)
	u.Data["json"] = Account{"CreateAccount", user, len(password), password}
	u.ServeJSON()
}

// @router /removeaccount [get]
func (u *ChainController) RemoveAccount() {
	fmt.Println("ChainController.RemoveAccount")
	password := u.GetString("password")
	user := u.GetString("user")
	fmt.Printf("password=%s, user=%s\n", password, user)
	u.Data["json"] = Account{"RemoveAccount", user, len(password), password}
	u.ServeJSON()
}

// @router /sendvote [post]
func (u *ChainController) SendVote() {
	fmt.Println("ChainController.SendVote")
	var tx TX
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &tx)
	if err != nil {
		fmt.Printf("bad RequestBody, %s\n", string(u.Ctx.Input.RequestBody))
	} else {
		u.Data["json"] = tx
	}
	u.ServeJSON()
}

// @router /sendtransaction [post]
func (u *ChainController) SendTransaction() {
	fmt.Println("ChainController.SendTransaction")
	var vt VT
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &vt)
	if err != nil {
		fmt.Printf("bad RequestBody, %s\n", string(u.Ctx.Input.RequestBody))
	} else {
		u.Data["json"] = vt
	}
	u.ServeJSON()
}
