package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"db"
	. "util"
)

//获取所有的商家信息
func GetAllCeller(session *JsNet.StSession) {
	list, err := GetAllSeller()
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	type Seller struct {
		UID    string
		OpenId string
		Name   string //姓名
		Mobile string //手机号
		Addr   string //地址
		Type   string //客户类型
		AdvNum int    //广告数量
	}

	data := []*Seller{}
	for _, v := range list {
		info := &Seller{}
		user, err := GetUserInfo(v)
		if err == nil {
			info.UID = user.UID
			info.OpenId = user.OpenId
			info.Name = user.Name
			info.Mobile = user.Mobile
			info.Addr = user.Addr
			info.Type = user.Type
		}
		AC, err := GetUserActivity(v)
		if err == nil {
			info.AdvNum = len(AC.Send)
		}
		data = append(data, info)
	}

	Forward(session, "0", data)
}

///提升用户为商家
func PromoteToSeller(UID string) error {
	data := &User{}
	if err := db.WriteLock(constant.Hash_User, UID, data); err != nil {
		return err
	}
	data.Role = "1"
	if err := db.WriteBack(constant.Hash_User, UID, data); err != nil {
		return err
	}
	return appendToGlobalSeller(UID)
}

///添加全局的商家列表
func appendToGlobalSeller(UID string) error {
	all := []string{}
	err := db.WriteLock(constant.Hash_User, constant.KEY_Global_Seller, &all)
	AppendUniqueString(&all, UID)
	if err != nil {
		return db.DirectWrite(constant.Hash_User, constant.KEY_Global_Seller, &all)
	}
	return db.WriteBack(constant.Hash_User, constant.KEY_Global_Seller, &all)
}

//获取所有的商家id
func GetAllSeller() ([]string, error) {
	all := []string{}
	err := db.ShareLock(constant.Hash_User, constant.KEY_Global_Seller, &all)
	return all, err
}

func RelationUserToDistribution(session *JsNet.StSession) {
	type INFO struct {
		UID    string //用户id
		Name   string //用户姓名
		CID    string //公司id
		Mobile string //手机号
		Code   string //短息验证码
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" || st.Mobile == "" || st.Code == "" || st.CID == "" {
		ForwardEx(session, "1", nil, "RelationUserToDistribution param failed , UID=%s,Name=%s,CID=%s,Mobile=%s,Code=%s\n",
			st.UID, st.Name, st.CID, st.Mobile, st.Code)
		return
	}

	if !VerifySmsCode(st.Mobile, st.Code) {
		ForwardEx(session, "1", nil, "短信校验失败\n")
		return
	}

	user := &User{}
	err := db.WriteLock(constant.Hash_User, st.UID, user)
	defer db.WriteBack(constant.Hash_User, st.UID, user)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	dis, err := getCompanyDistribution(st.CID)
	did := []string{}
	exist := false
	if err != nil {
		Warn("该公司没有网点,err:" + err.Error())
	} else {
		for _, v := range dis.DistributionID {
			d, err := GetDistributionInfo(v)
			if err == nil {
				if d.DMobile == st.Mobile {
					did = append(did, v)
					exist = true
				}
			}
		}
	}

	if user.Mobile == st.Mobile {
		exist = true
	}

	if !exist {
		Forward(session, "1", nil)
		return
	} else {
		ok := false
		for i, v := range user.IsManager {
			if v.CID == st.CID {
				if len(did) > 0 {
					for _, v1 := range v.DID {
						for _, v2 := range did {
							if v1 != v2 {
								user.IsManager[i].DID = append(user.IsManager[i].DID, v2)
							}
						}
					}
				}
				ok = true
				break
			}
		}
		if !ok {
			ma := ManagerID{CID: st.CID}
			if len(did) > 0 {
				ma.DID = append(ma.DID, did...)
			}
			user.IsManager = append(user.IsManager, ma)
		}
	}
	Forward(session, "0", nil)
}
