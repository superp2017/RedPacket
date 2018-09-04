package common

import (
	."JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"db"
	"ider"
	."util"
)

type AnswerInfo struct {
	Content string //答案
	Checked bool   //是否选中
}

type ActivityQuestion struct {
	QID        string       //问题id
	Question   string       //问题
	Answer     string       //答案
	AnswerList []AnswerInfo //答案列表
	Date       string       //创建日期
}

func NewQuestion(session *JsNet.StSession) {
	que := &ActivityQuestion{}
	if err := session.GetPara(que); err != nil {
		ForwardEx(session, "1", nil, "NewQuestion GetPara:%s\n", err.Error())
		return
	}
	if que.Answer == "" || que.Question == "" || len(que.AnswerList) == 0 {
		ForwardEx(session, "1", nil, "NewQuestion failed,Question=%s,Answer=%s,len(que.AnswerList)=%d\n", que.Question, que.Answer, len(que.AnswerList))
		return
	}
	que.QID = ider.GenID()
	que.Date = CurTime()

	if err := db.DirectWrite(constant.Hash_ActivityQuestion, que.QID, que); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", que)
}
//创建问题
func CreatQuetion(Q, A string, list []AnswerInfo) (*ActivityQuestion, error) {
	que := &ActivityQuestion{}
	if Q == "" || A == "" || len(list) == 0 {
		return que, ErrorLog("CreatQuetion,Param is empty,Que:%s,Answer:%s,Len:%d", Q, A, len(list))
	}
	que.QID = ider.GenID()
	que.Date = CurTime()
	que.Question = Q
	que.Answer = A
	que.AnswerList = list
	err := db.DirectWrite(constant.Hash_ActivityQuestion, que.QID, que)
	return que, err
}

//查询问题
func QueryQuertion(session *JsNet.StSession) {
	type INFO struct {
		QID string //问题id
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	que, err := GetQuestion(st.QID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", que)
}

///活动绑定问题
func ActivityBindQuestion(session *JsNet.StSession) {
	type Settings struct {
		ActivityID string   //活动Id
		Question   []string //问题id列表
	}
	st := &Settings{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "ActivityBindQuestion GetPara:%s\n", err.Error())
		return
	}
	if st.ActivityID == "" || len(st.Question) == 0 {
		ForwardEx(session, "1", nil, "ActivityBindQuestion Param is empty ActivityID:%s,Question=%v\n", st.ActivityID, st.Question)
		return
	}
	data := &Activity{}
	if err := db.WriteLock(constant.Hash_Activity, st.ActivityID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

///获取活动关联的问题
func GetActivityQuestion(session *JsNet.StSession) {
	type INFO struct {
		ActivityID string
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	ac, err := GetActivityInfo(st.ActivityID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	data := []*ActivityQuestion{}
	for _, v := range ac.QuestionId {
		que, er := GetQuestion(v)
		if er == nil {
			data = append(data, que)
		}
	}
	Forward(session, "0", data)
}

//修改问题
func ModifyQuestion(session *JsNet.StSession) {
	type INFO struct {
		QID        string       //问题id
		Question   string       //问题
		Answer     string       //答案
		AnswerList []AnswerInfo //答案列表
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.QID == "" || st.Answer == "" || st.Question == "" || len(st.AnswerList) == 0 {
		ForwardEx(session, "1", nil,
			"ModifyQuestion failed,QID=%s,Question=%s,Answer=%s,AnswerList len=%d\n",
			st.QID, st.Question, st.Answer, len(st.AnswerList))
		return
	}
	que := &ActivityQuestion{}
	if err := db.WriteLock(constant.Hash_ActivityQuestion, que.QID, que); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	que.Question = st.Question
	que.Answer = st.Answer
	que.AnswerList = st.AnswerList

	if err := db.WriteBack(constant.Hash_ActivityQuestion, que.QID, que); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", que)
}

func GetQuestion(QID string) (*ActivityQuestion, error) {
	if QID == "" {
		return nil, ErrorLog("GetQuestion failed,QID is empty !\n")
	}
	que := &ActivityQuestion{}
	err := db.ShareLock(constant.Hash_ActivityQuestion, QID, que)
	return que, err
}
