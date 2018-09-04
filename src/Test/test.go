package main

import (
	"JsLib/JsDispatcher"
	"JsLib/JsExit"
	"JsLib/JsNet"
	"common"
	"constant"
	"db"
	. "util"
)

////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
/////////////////////////这个类用于紧急情况维护,不到万不得已不要调用////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////

func exit() int {
	JsDispatcher.Close()
	return 0
}

func main() {
	JsExit.RegisterExitCb(exit)
	JsNet.AppConf("./conf/app.conf")
	initRouter()
	JsDispatcher.Run()
}

func initRouter() {
	JsNet.Http("/UpdataCompanyAccount", UpdataCompanyAccount) ///重新建立所有商家的账户表
	JsNet.Http("/GetKeys", GetKeys)                           ///获取某一个表的所有的keys
	JsNet.Http("/GetWarningUserMoney", GetWarningUserMoney)   //检查用户余额或者收到钱的异常用户
	JsNet.Http("/reMapUnionID2UID", reMapUnionID2UID)         //检查用户余额或者收到钱的异常用户
}

///重新建立所有商家的账户表
func UpdataCompanyAccount(session *JsNet.StSession) {
	list, err := common.GetGlobalCompany()
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	for _, v := range list {
		com, err := common.GetCompanyInfo(v)
		if err == nil {
			common.AddCompanyAccount(com.UserName, com.Password, com.FullName, com.CID)
		}
	}
	Forward(session, "0", nil)
}

///获取某一个表的所有的keys
func GetKeys(session *JsNet.StSession) {
	type INFO struct {
		TABLE string
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.TABLE == "" {
		ForwardEx(session, "1", nil, "GetKeys param is empty\n")
		return
	}
	Forward(session, "0", db.GetKeys(st.TABLE))
}

func reMapUnionID2UID(session *JsNet.StSession) {
	type INFO struct {
		UIDs  []string
		IsAll bool
	}

	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if !st.IsAll && len(st.UIDs) == 0 {
		ForwardEx(session, "1", nil, "ResetUserMoney failed,IsAll:%v,IsAll:%v\n", st.IsAll, st.UIDs)
		return
	}

	if st.IsAll {
		st.UIDs = db.GetKeys(constant.Hash_User)
	}

	for _, v := range st.UIDs {
		user := &common.User{}
		if err := db.ShareLock(constant.Hash_User, v, user); err == nil {
			go common.OpenidMapUID(user)
			go common.UnionidMapUID(user)
		}
	}
	Forward(session, "0", nil)
}

///获取用户金额异常的用户
func GetWarningUserMoney(session *JsNet.StSession) {
	type INFO struct {
		Money int
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	ids := db.GetKeys(constant.Hash_User)
	data := []string{}
	for _, v := range ids {
		user := &common.User{}
		if err := db.ShareLock(constant.Hash_User, v, user); err == nil {
			if user.Blance > st.Money || user.RecevieMoney > st.Money {
				data = append(data, user.UID)
			}
		}
	}
	Forward(session, "0", data)
}

///重置多个或者全部用户的账户信息
func ResetUserMoney(session *JsNet.StSession) {
	type INFO struct {
		UIDs  []string
		IsAll bool
	}

	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if !st.IsAll && len(st.UIDs) == 0 {
		ForwardEx(session, "1", nil, "ResetUserMoney failed,IsAll:%v,IsAll:%v\n", st.IsAll, st.UIDs)
		return
	}

	if st.IsAll {
		st.UIDs = db.GetKeys(constant.Hash_User)
	}

	for _, v := range st.UIDs {
		user := &common.User{}
		if err := db.WriteLock(constant.Hash_User, v, user); err == nil {
			user.Blance = 0
			user.RecevieMoney = 0
			user.Recharge = 0
			db.WriteBack(constant.Hash_User, v, user)
		}
	}
	Forward(session, "0", nil)
}

///重置UserActivity部分或者全部信息
func ResetUserActity(session *JsNet.StSession) {
	type INFO struct {
		UIDs      []string ///用户id
		IsAll     bool     //是否全部用户
		IsSent    bool     //是否重置Sent
		IsRecevie bool     //是否重置Recevie
		IsRequire bool     //是否重置需求
		IsMoney   bool     //是否充值Money
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if !st.IsAll && len(st.UIDs) == 0 {
		ForwardEx(session, "1", nil, "ResetUserActivity param failed,UIDs=%v,IsAll%v\n", st.UIDs, st.IsAll)
		return
	}
	if st.IsAll {
		st.UIDs = db.GetKeys(constant.Hash_UserActivity)
	}
	for _, v := range st.UIDs {
		userAC := &common.UserActivity{}
		if err := db.WriteLock(constant.Hash_UserActivity, v, userAC); err == nil {
			if st.IsSent {
				userAC.Send = []common.AcADP{}
				userAC.Send_Detail = []*common.Activity{}
			}
			if st.IsRecevie {
				userAC.Receive = []common.AcADP{}
				userAC.Receive_Detail = []*common.Activity{}
			}
			if st.IsRequire {
				userAC.Requirement = []string{}
				userAC.Re_Detail = []*common.Requirement{}
			}
			if st.IsMoney {
				userAC.SendMoney = 0
				userAC.RecevieMoney = 0
			}
			db.WriteBack(constant.Hash_UserActivity, v, userAC)
		}
	}
	Forward(session, "0", nil)
}
