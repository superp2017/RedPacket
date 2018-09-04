package common

import (
	"JsLib/JsNet"
	"constant"
	"db"
	"ider"
	. "util"
)

type feedBack struct {
	FID     string //反馈ID
	ID      string //用户或者商家的id
	Title   string //标题
	Content string //正文
	Name    string //姓名
	Mobile  string //手机号
	Direct  int    //来源0:用户 1:商家
	Status  int    //状态 0:新建 1:已经受理
	Date    string //日期
}

func newFeedBack(session *JsNet.StSession) {
	st := &feedBack{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.Content == "" {
		ForwardEx(session, "1", nil, "newFeedBack ,Content isEmty,data:%v", st)
		return
	}
	st.FID = ider.GenID()
	st.Date = CurTime()
	st.Status = 0
	if err := db.DirectWrite(constant.Hash_FeedBack, st.FID, st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	go addGlobalFeedBack(st.FID)
	Forward(session, "0", st)
}

func queryFeedBack(session *JsNet.StSession) {
	type INFO struct {
		FID string
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.FID == "" {
		ForwardEx(session, "1", nil, "queryFeedBack FID isEmpty\n")
		return
	}
	data, err := getFeedBack(st.FID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func getGlobalFeedBack(session *JsNet.StSession) {
	data := []*feedBack{}
	ids := []string{}
	db.ShareLock(constant.Hash_FeedBack, constant.KEY_GlobalFeedBack, &ids)
	for _, v := range ids {
		if d, e := getFeedBack(v); e != nil {
			data = append(data, d)
		}
	}
	Forward(session, "0", data)
}

func getFeedBack(FID string) (*feedBack, error) {
	st := &feedBack{}
	err := db.ShareLock(constant.Hash_FeedBack, FID, st)
	return st, err
}

func addGlobalFeedBack(FID string) error {
	data := []string{}
	err := db.WriteLock(constant.Hash_FeedBack, constant.KEY_GlobalFeedBack, &data)
	AppendUniqueString(&data, FID)
	if err != nil {
		return db.DirectWrite(constant.Hash_FeedBack, constant.KEY_GlobalFeedBack, &data)
	}
	return db.WriteBack(constant.Hash_FeedBack, constant.KEY_GlobalFeedBack, &data)
}
