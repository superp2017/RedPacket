package common

import (
	"JsLib/JsNet"
	"constant"
	"db"
	"ider"
	. "util"
)

type AgentAccount struct {
	UserName   string //登陆用户名
	AgentID    string //公司ID
	Password   string //登陆密码
	EntityTime string //创建日期
}

//代理登陆
func AgentLogIn(session *JsNet.StSession) {
	type RD_Login struct {
		UserName string //用户名
		PassWord string //用户密码
	}
	st := RD_Login{}
	if err := session.GetPara(&st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	agentAccount := &AgentAccount{}
	if err := db.ShareLock(constant.Hash_Agent_Account, st.UserName, agentAccount); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	data := &User{}
	if err := db.ShareLock(constant.Hash_User, agentAccount.AgentID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if data.UserName != st.UserName || data.PassWord != st.PassWord {
		ForwardEx(session, "1", nil, "The User Name or the Password is not correct")
		return
	}
	Forward(session, "0", data)
}

//绑定UID
func BindUser(session *JsNet.StSession) {
	type BindInfo struct {
		AgentUID    string //代理UID
		FollowerUID string //被绑定的UID
	}
	st := BindInfo{}
	if err := session.GetPara(&st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	agent := &User{}
	if err := db.WriteLock(constant.Hash_User, st.AgentUID, agent); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	follower := &User{}
	if err := db.WriteLock(constant.Hash_User, st.FollowerUID, follower); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		db.WriteBack(constant.Hash_User, st.AgentUID, agent)
		return
	}

	if agent.Role != "2" {
		ForwardEx(session, "1", nil, "user is not the agent")
		db.WriteBack(constant.Hash_User, st.FollowerUID, follower)
		db.WriteBack(constant.Hash_User, st.AgentUID, agent)
		return
	}

	if follower.Role != "2" {
		ForwardEx(session, "1", nil, "follower is  the agent")
		db.WriteBack(constant.Hash_User, st.FollowerUID, follower)
		db.WriteBack(constant.Hash_User, st.AgentUID, agent)
		return
	}

	follower.ParentAgentID = st.AgentUID

	agentFollower := AgentFollower{}
	agentFollower.EntityTime = CurTime()
	agentFollower.Name = follower.Nickname
	agentFollower.UID = follower.UID
	agentFollower.Address = follower.Addr

	appendnow := true
	for _, v := range agent.LsAgentFollower {
		if v.UID == follower.UID {
			appendnow = false
			break
		}
	}
	if appendnow {
		agent.LsAgentFollower = append(agent.LsAgentFollower, agentFollower)
	}

	db.WriteBack(constant.Hash_User, st.FollowerUID, follower)
	db.WriteBack(constant.Hash_User, st.AgentUID, agent)
	Forward(session, "0", nil)
}

/*
 注册一个公司账号
*/
func AgentNewCompany(session *JsNet.StSession) {
	st := &Company{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.FullName == "" {
		ForwardEx(session, "1", nil, "New Company failed with FullName=%s and UserName=%s,Password=%s\n", st.FullName, st.UserName, st.Password)
		return
	}

	if st.AgentID == "" {
		ForwardEx(session, "1", nil, "New Company failed with AgentID=%s\n", st.AgentID)
		return
	}
	st.CID = ider.GenID()
	st.UserName = st.FullName
	st.Password = st.CID
	st.IsAgentLaunched = true

	// if CheckCompanyAccount(st.UserName) {
	// 	ForwardEx(session, "1", nil, " 用户名已经存在,UserName=%s", st.UserName)
	// 	return
	// }

	st.EntityTime = CurTime()
	if err := db.DirectWrite(constant.Hash_Company, st.CID, st); err != nil {
		ForwardEx(session, "1", nil, "New Company DirectWrite :"+err.Error())
		return
	}
	cudb := &CUser{}
	cudb.CID = st.CID
	cudb.HmCUser = make(map[string][]ShortInfo)
	if err := db.DirectWrite(constant.Hash_CRecordUser, cudb.CID, cudb); err != nil {
		ForwardEx(session, "1", nil, "New Company DirectWrite :"+err.Error())
		return
	}
	///添加公司账号
	if err := AddCompanyAccount(st.UserName, st.Password, st.FullName, st.CID); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	//公司信息添加到代理
	companShortInfo := CompanyShortInfo{}
	companShortInfo.CompanyID = st.CID
	companShortInfo.CompanyName = st.FullName
	companShortInfo.EntityTime = st.EntityTime
	user := &User{}
	err := db.WriteLock(constant.Hash_User, st.AgentID, user)
	if err != nil {
		Forward(session, "1", st)
		return
	}

	user.LsCompanyShortInfo = append(user.LsCompanyShortInfo, companShortInfo)
	db.WriteBack(constant.Hash_User, st.AgentID, user)

	user.LsCompanyShortInfo = append(user.LsCompanyShortInfo, companShortInfo)
	db.WriteLock(constant.Hash_User, st.AgentID, user)
	Forward(session, "0", st)
}

func NewAgent(session *JsNet.StSession) {
	//pull the user out
	type UserInfo struct {
		UID      string
		UserName string
		PassWord string
	}
	st := &UserInfo{}
	user := &User{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if st.UID == "" || st.UserName == "" || st.PassWord == "" {
		ForwardEx(session, "1", nil, "The User's UID=%s,UserName=%s,Password=%s wrong\n", st.UID, st.UserName, st.PassWord)
		return
	}
	ok := IsAgentAccountExit(st.UserName)
	if ok {
		ForwardEx(session, "1", nil, "用户名存在，请更换其他的用户名\n", st.UserName)
		return
	}
	err := db.WriteLock(constant.Hash_User, st.UID, user)
	if err != nil {
		Forward(session, "1", st)
		return
	}
	if user.Role == "2" {
		db.WriteBack(constant.Hash_User, st.UID, user)
		ForwardEx(session, "1", nil, "The User is already the agent at the time =%s\n", user.AgentTime)
		return
	}
	user.Role = "2"
	user.UserName = st.UserName
	user.PassWord = st.PassWord
	user.AgentTime = CurTime()
	db.WriteBack(constant.Hash_User, st.UID, user)

	agentAccont := &AgentAccount{}
	agentAccont.UserName = st.UserName
	agentAccont.EntityTime = CurTime()
	agentAccont.Password = st.PassWord
	agentAccont.AgentID = st.UID
	err = db.DirectWrite(constant.Hash_Agent_Account, agentAccont.UserName, agentAccont)
	if err != nil {
		Forward(session, "1", nil)
		return
	}
	Forward(session, "0", st)
}

///检查用户名是否重复
func IsAgentAccountExit(UserName string) bool {
	data := &AgentAccount{}
	if err := db.ShareLock(constant.Hash_Agent_Account, UserName, data); err != nil {
		return false
	}
	return true
}
