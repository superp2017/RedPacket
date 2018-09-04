package common

import (
	"JsLib/JsNet"
	"ider"
	"sort"
	"sync"
	. "util"
)

var LOGIN_TOKEN map[string]string = make(map[string]string) //全局token
var gl_token_mutex sync.Mutex

type RDPRecord []ActivityRDP

func (list RDPRecord) Len() int {
	return len(list)
}

func (list RDPRecord) Less(i, j int) bool {
	return list[i].TimeStap < list[j].TimeStap
}

func (list RDPRecord) Swap(i, j int) {
	var temp ActivityRDP = list[i]
	list[i] = list[j]
	list[j] = temp
}

//查询后台活动的统计信息
func QueryActivityStatistics(session *JsNet.StSession) {
	type INFO struct {
		SortType int //排序方式  0：小时 1：天 2 ：月
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	type Ret struct {
		Data  map[string]RDPRecord ///key是刻度  RDPRecond 是[]ActivityRDP,所有领取记录
		Scale []string             //刻度数组
	}
	ret := Ret{}
	ret.Data = make(map[string]RDPRecord)
	ret.Scale = []string{}

	list := getAllActivityID()
	if len(list) == 0 {
		Forward(session, "0", ret)
		return
	}
	data := RDPRecord{}
	for _, v := range list {
		rdp, err := getActivityRDP(v)
		if err == nil {
			data = append(data, rdp.RDP...)
		}
	}

	if len(data) == 0 {
		Forward(session, "0", ret)
		return
	}
	sort.Sort(data)

	for _, v := range data {
		str := ""
		if st.SortType == 0 {
			t, e := GetTimeFormString(v.Date)
			if e != nil {
				continue
			}
			str = GetYMDH_CH(t)
		} else if st.SortType == 1 {
			t, e := GetTimeFormString(v.Date)
			if e != nil {
				continue
			}
			str = GetYMD_CH(t)
		} else if st.SortType == 2 {
			t, e := GetTimeFormString(v.Date)
			if e != nil {
				continue
			}
			str = GetYM_CH(t)
		}
		if l, ok := ret.Data[str]; ok {
			l = append(l, v)
			ret.Data[str] = l
		} else {
			L := RDPRecord{v}
			ret.Data[str] = L
		}
		AppendUniqueString(&ret.Scale, str)
	}

	Forward(session, "0", ret)
}

func BackTrans(session *JsNet.StSession) {
	type info struct {
		UID   string
		Money int
		Token string
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if ok := checkToken(session.RemoteAddr(), st.Token); !ok {
		ForwardEx(session, "1", nil, "BackTrans,Token 校验失败\n")
		return
	}

	user, err := GetUserInfo(st.UID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data := &ST_Transfer{
		OpenId:     user.OpenId,
		UserName:   user.Name,
		UserHeader: user.HeadImageURL,
		TimeStamp:  CurStamp(),
		LDate:      CurTime(),
		Desc:       "传单侠奖励红包",
		Amount:     st.Money,
	}
	cb, err := direct_transfer(data)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", cb)
}

func BackPayRequirment(session *JsNet.StSession) {
	type info struct {
		RequirementID string //需求id
		Token         string //令牌
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if ok := checkToken(session.RemoteAddr(), st.Token); !ok {
		ForwardEx(session, "1", nil, "BackPayRequirment,Token 校验失败\n")
		return
	}

	if st.RequirementID == "" {
		ForwardEx(session, "1", nil, "BackPayRequirment,RequirementID 为空\n")
		return
	}

	if err := PayRequirement(st.RequirementID); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", nil)
}

func BackPayActity(session *JsNet.StSession) {
	type info struct {
		ActivityID string //活动id
		Token      string //令牌
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if ok := checkToken(session.RemoteAddr(), st.Token); !ok {
		ForwardEx(session, "1", nil, "BackPayActity,Token 校验失败\n")
		return
	}

	if st.ActivityID == "" {
		ForwardEx(session, "1", nil, "BackPayActity,ActivityID 为空\n")
		return
	}

	if err := PayActivity(st.ActivityID); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", nil)
}

//检验token
func checkToken(ip, token string) bool {
	if t, ok := LOGIN_TOKEN[ip]; ok {
		if t == token {
			return true
		}
	}
	return false
}

//更新token
func MapToken(IP string) string {
	id := ider.GenID()
	gl_token_mutex.Lock()
	LOGIN_TOKEN[IP] = id
	gl_token_mutex.Unlock()
	return id
}

//清除token
func ClearToken(IP string) {
	gl_token_mutex.Lock()
	LOGIN_TOKEN[IP] = ""
	gl_token_mutex.Unlock()
}
