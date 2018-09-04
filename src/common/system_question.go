package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"db"
	"errors"
	"ider"
	. "util"
)

type SystemQuestionId struct {
	QuestionLevel            string         //问题id
	HmQuestion               map[int]string //系统问题Id
	SystemQuestionTotalCount int
	SystemQuestionValidCount int
}
type SystemQuestion struct {
	Question   string       //问题
	EntityTime       string //创建日期
	AnswerList []string     //答案,注意这个是用户的选择,包括
	ChoiceList []AnswerInfo //备选答案列表
	QID        string       //问题id
	CID        string       //公司ID
	CName      string       //公司名称
	ChoiceTips  string      //选项提示
	ChoiceAnswer string     //答案
	QuestionType string     //"Choice";"InputAnswer","ChoiseInput"
	ChoiseType   string     //"Combobox","List"
	InputAnswerTips []string  //答案的备选选项
	QuestionDex int          //问题的标签,即给客户显示的顺序，是先这个
}
type UserSystemQuestion struct {
	UID                      string
	AnswerCountDex           int
	AnswerTotalCount         int
	LsSystemQuestionAnswered []string //已经回答过的问题列表
	LsSystemQuestionWaitting []string //目前需要回答的系统问题
}

//
/*
 新建一系统问题
*/
func NewSystemQuetion(session *JsNet.StSession) {
	//Create the system question
	st := &SystemQuestion{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.Question == "" || st.QuestionType == "" || st.QuestionDex == 0 {

		ForwardEx(session, "1", nil, "CreatQuetion,Param is empty,Que:%s,QuestionType:%s", st.Question, st.QuestionType)
		return
	}

	if st.QuestionType == constant.Question_ChoiseInput || st.QuestionType == constant.Question_Choise {
		if len(st.ChoiceList) == 0 {
			ForwardEx(session, "1", nil, "Choise list is 0")
			return
		}
	}
	st.QID = ider.GenID()
	st.EntityTime = CurTime()
	//put the system question to the db
	if err := db.DirectWrite(constant.Hash_SystemQuestion, st.QID, st); err != nil {
		ForwardEx(session, "1", nil, "NewRequirement DirectWrite :"+err.Error())
		return
	}
	// get the systemquestion id out
	questionIDEntity := &SystemQuestionId{}
	err := db.WriteLock(constant.Hash_SystemQuestionID, constant.KEY_SYSQuestion_Normal, questionIDEntity)
	if err != nil {
		questionIDEntity.QuestionLevel = constant.KEY_SYSQuestion_Normal
		questionIDEntity.HmQuestion = make(map[int]string)
		questionIDEntity.HmQuestion[st.QuestionDex] = st.QID
		questionIDEntity.SystemQuestionTotalCount = 1
		questionIDEntity.SystemQuestionValidCount=1
		db.DirectWrite(constant.Hash_SystemQuestionID, constant.KEY_SYSQuestion_Normal, questionIDEntity)
	} else {
		questionIDEntity.HmQuestion[st.QuestionDex] = st.QID
		questionIDEntity.SystemQuestionTotalCount = questionIDEntity.SystemQuestionTotalCount + 1
		lcount:=0
		for _,_= range questionIDEntity.HmQuestion {
			lcount=lcount+1
		}
		questionIDEntity.SystemQuestionValidCount=lcount
		db.WriteBack(constant.Hash_SystemQuestionID, constant.KEY_SYSQuestion_Normal, questionIDEntity)
	}
	Forward(session, "0", st)
}


//新建公司的一个系统问题

func NewActivitySystemQuestion(session *JsNet.StSession) {
	//Create the system question
	type AllID struct{
		LsAllID []string
	}
	st := &SystemQuestion{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.Question == "" || st.QuestionType == "" || st.QuestionDex == 0 {
		ForwardEx(session, "1", nil, "CreatQuetion,Param is empty,Que:%s,QuestionType:%s", st.Question, st.QuestionType)
		return
	}
	if st.QuestionType == constant.Question_ChoiseInput || st.QuestionType == constant.Question_Choise {
		if len(st.ChoiceList) == 0 {
			ForwardEx(session, "1", nil, "Choise list is 0")
			return
		}
	}
	st.QID = ider.GenID()
	st.EntityTime = CurTime()
	//put the activity system question to the db
	if err := db.DirectWrite(constant.Hash_ActivitySQuestion,st.QID,st); err != nil {
		ForwardEx(session, "1", nil, "NewRequirement DirectWrite :"+err.Error())
		return
	}

	err := AppendSAQuestionToCompany(st.QID, st.CID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	//put the systemquestion to the allID
	allID:=&AllID{}
	if err := db.WriteLock(constant.Hash_ActivitySQuestion,"AllID",allID); err != nil {
		allID.LsAllID=[]string{}
		allID.LsAllID=append(allID.LsAllID,st.QID)
		db.DirectWrite(constant.Hash_ActivitySQuestion,"AllID",allID)
	}else{
		allID.LsAllID=append(allID.LsAllID,st.QID)
		db.WriteBack(constant.Hash_ActivitySQuestion,"AllID",allID)
	}
	Forward(session, "0", st)
}






func NewActivitySystemQuestionLocal(st *SystemQuestion)(e error) {
	//Create the system question
	if st.Question == "" || st.QuestionType == "" || st.QuestionDex == 0 {
		return errors.New("CreatQuetion,Param is empty")
	}
	if st.QuestionType == constant.Question_ChoiseInput || st.QuestionType == constant.Question_Choise {
		if len(st.ChoiceList) == 0 {
			return errors.New("Choise list is 0")
		}
	}
	st.QID = ider.GenID()
	st.EntityTime = CurTime()
	//put the activity system question to the db
	if err := db.DirectWrite(constant.Hash_ActivitySQuestion,st.QID,st); err != nil {
		return err
	}
	return nil
}





/*
*/

// ///创建系统问题
// func CreatSystemQuetionBk(question, questionType string, lsChoise []AnswerInfo) (*SystemQuestion, error) {
// 	que := &SystemQuestion{}
// 	if question == "" || questionType == "" {
// 		return que, ErrorLog("CreatQuetion,Param is empty,Que:%s,QuestionType:%s,Len:%d", question, questionType, len(lsChoise))
// 	}

// 	if questionType == constant.Question_ChoiseInput || questionType == constant.Question_Choise {
// 		if len(lsChoise) == 0 {
// 			return que, ErrorLog("CreatQuetion,Param is empty,Que:%s,QuestionType:%s,Len:%d", question, questionType, len(lsChoise))
// 		}
// 	}
// 	que.QID = ider.GenID()
// 	que.EntityTime = CurTime()
// 	que.Question = question
// 	que.QuestionType = questionType
// 	que.ChoiceList = lsChoise
// 	err := db.DirectWrite(constant.Hash_SystemQuestion, que.QID, que)
// 	return que, err
// }

///创建系统问题
func CreatSystemQuetionBk(question, questionType string, lsChoise []AnswerInfo) (*SystemQuestion, error) {
	que := &SystemQuestion{}
	if question == "" || questionType == "" {
		return que, ErrorLog("CreatQuetion,Param is empty,Que:%s,QuestionType:%s,Len:%d", question, questionType, len(lsChoise))
	}

	if questionType == constant.Question_ChoiseInput || questionType == constant.Question_Choise {
		if len(lsChoise) == 0 {
			return que, ErrorLog("CreatQuetion,Param is empty,Que:%s,QuestionType:%s,Len:%d", question, questionType, len(lsChoise))
		}
	}
	que.QID = ider.GenID()
	que.EntityTime = CurTime()
	que.Question = question
	que.QuestionType = questionType
	que.ChoiceList = lsChoise
	err := db.DirectWrite(constant.Hash_SystemQuestion, que.QID, que)
	return que, err
}

//查询问题
func QuerySystemQuertion(session *JsNet.StSession) {
	type INFO struct {
		QID string //问题id
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	que, err := GetSystemQuestionInfo(st.QID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", que)
}


//查询问题
func QueryASystemQuertion(session *JsNet.StSession) {
	type INFO struct {
		QID string //问题id
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	que, err := GetASystemQuestionInfo(st.QID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", que)
}




//查询多个系统问题
func QueryMultiASystemQuertion(session *JsNet.StSession) {
	type INFO struct {
		LsASystemQuestion []string //问题id
	}
	LsQuestion:=[]*SystemQuestion{}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	for i:=0;i<len(st.LsASystemQuestion);i++{
		que, err := GetASystemQuestionInfo(st.LsASystemQuestion[i])
		if err==nil{
			LsQuestion=append(LsQuestion,que)
		}
	}
	
	Forward(session, "0", LsQuestion)
}


//查询所有公司问题
func QueryAllCompanySystemAQuestion(session *JsNet.StSession) {
	type AllID struct{
		LsAllID []string
	}
	allID:=&AllID{}
	LsQuestion:=[]*SystemQuestion{}
	if err := db.ShareLock(constant.Hash_ActivitySQuestion,"AllID",allID); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	for _,v:=range allID.LsAllID{
		que, err := GetASystemQuestionInfo(v)
		if err==nil{
			LsQuestion=append(LsQuestion,que)
		}
	}
	Forward(session, "0", LsQuestion)
	}
	



//查询所有系统问题
func QueryAllSystemQuertion(session *JsNet.StSession) {

	systemQuestionIdDB := &SystemQuestionId{}
	LsQuestion:=[]*SystemQuestion{}
	//drag out the system question ids
	err := db.ShareLock(constant.Hash_SystemQuestionID, constant.KEY_SYSQuestion_Normal, systemQuestionIdDB)
	if  err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	for _,v:=range systemQuestionIdDB.HmQuestion{
		que, err := GetSystemQuestionInfo(v)
		if err==nil{
			LsQuestion=append(LsQuestion,que)
		}
	}
	Forward(session, "0", LsQuestion)
}






///活动绑定系统问题
func ActivityBindSystemQuestion(ActivityID string, lsQuestion []string) (activity *Activity, e error) {
	data := &Activity{}
	if ActivityID == "" || len(lsQuestion) == 0 {
		return data, errors.New("Activity ID is none or lsQuestion length is 0")
	}
	if err := db.WriteLock(constant.Hash_Activity, ActivityID, data); err != nil {
		return data, errors.New("There is no such an activity in the database")
	}

	if len(data.SystemQuestionId) == 0 {
		data.SystemQuestionId = []string{}
	}

	for _, v := range lsQuestion {
		data.SystemQuestionId = append(data.SystemQuestionId, v)
	}

	if err := db.WriteBack(constant.Hash_Activity, ActivityID, data); err != nil {
		return data, errors.New("Write Activity back wrong")
	}
	return data, errors.New("Write Activity back wrong")
}

///获取活动关联的系统问题
func GetSystemQuestion(session *JsNet.StSession) {
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

	data := []*SystemQuestion{}
	for _, v := range ac.SystemQuestionId {
		que, er := GetSystemQuestionInfo(v)
		if er == nil {
			data = append(data, que)
		}
	}
	Forward(session, "0", data)
}

//删除活动的系统问题
func DeleteSystemQuestion(session *JsNet.StSession) {
	type INFO struct {
		QuestionID string //问题id
		ActivityID string //活动id
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.QuestionID == "" || st.ActivityID == "" {
		ForwardEx(session, "1", nil,
			"DeleteSystem question failed,QuestionID=%s,ActivityID=%s\n", st.QuestionID, st.ActivityID)
		return
	}
	activityDB := &Activity{}
	if err := db.WriteLock(constant.Hash_Activity, st.ActivityID, activityDB); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	DelExistString(&activityDB.SystemQuestionId, st.QuestionID)
	if err := db.WriteBack(constant.Hash_Activity, st.ActivityID, activityDB); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", activityDB)
}







//获取系统问题
func GetSystemQuestionInfo(QID string) (*SystemQuestion, error) {
	if QID == "" {
		return nil, ErrorLog("GetQuestion failed,QID is empty !\n")
	}
	que := &SystemQuestion{}
	err := db.ShareLock(constant.Hash_SystemQuestion, QID, que)
	return que, err
}


//获取活动的系统问题
func GetASystemQuestionInfo(QID string) (*SystemQuestion, error) {
	if QID == "" {
		return nil, ErrorLog("GetQuestion failed,QID is empty !\n")
	}
	que := &SystemQuestion{}
	err := db.ShareLock(constant.Hash_ActivitySQuestion, QID, que)
	return que, err
}


//系统自动分配用户若干个系统问题，这个在用户问题回答完成后自动更新
func UpdateUserSystemQuestion(userQuestion *UserSystemQuestion, initialNum int) (e error) {
	userQuestion.LsSystemQuestionWaitting = []string{}
	anseweredQuestionCount := userQuestion.AnswerCountDex
	systemQuestionIdDB := &SystemQuestionId{}
	//drag out the system question ids
	err := db.ShareLock(constant.Hash_SystemQuestionID, constant.KEY_SYSQuestion_Normal, systemQuestionIdDB)
	if err != nil {
		return ErrorLog("Pull the System Question ID fail")
	}
	initialDex := 1
	for i := 0; i < initialNum; i++ {
		if _, ok := systemQuestionIdDB.HmQuestion[anseweredQuestionCount+i]; ok {
			userQuestion.LsSystemQuestionWaitting = append(userQuestion.LsSystemQuestionWaitting, systemQuestionIdDB.HmQuestion[anseweredQuestionCount+i])
		} else {
			if _, okd := systemQuestionIdDB.HmQuestion[initialDex]; okd {
				userQuestion.LsSystemQuestionWaitting = append(userQuestion.LsSystemQuestionWaitting, systemQuestionIdDB.HmQuestion[initialDex])
				initialDex = initialDex + 1
			} else {
				continue
			}

		}
	}
	
	return nil

}

//获取网络层系统问题

//查询问题
func GetUserNetSystemQuestion(session *JsNet.StSession) {
	type INFO struct {
		UID string //问题id
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	lsQuestion, err := GetUserWaitingSystemQuestion(st.UID)

	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", lsQuestion)
}

//获取用户的系统问题

func GetUserWaitingSystemQuestion(UID string) (lsQuestion []*SystemQuestion, e error) {

	//load system question 
	user := &UserSystemQuestion{}
	err := db.ShareLock(constant.Hash_UserSystemQuestion, UID, user)
	if err != nil {
		user.UID=UID
		user.AnswerCountDex=0
		user.AnswerTotalCount=0
		UpdateUserSystemQuestion(user , 5)
		db.DirectWrite(constant.Hash_UserSystemQuestion, user.UID, user)
	}
	lsSystemQuestion := []*SystemQuestion{}
	
	for i := 0; i < len(user.LsSystemQuestionWaitting); i++ {
		systemQuestion, err := GetSystemQuestionInfo(user.LsSystemQuestionWaitting[i])
		if err == nil {
			lsSystemQuestion = append(lsSystemQuestion, systemQuestion)

		}
	}
	return lsSystemQuestion, nil
}

// AnswerActivitySQuestion(UID,AID string, LsSystemQuestion []SystemQuestion)

// AnswerSystemQuestion(UID string, LsSystemQuestion []SystemQuestion)


//用户回答系统问题，这个在用户领取红包时候内部调用，把系统问题传过来
func AnswerSystemQuestion(UID string, LsSystemQuestion []SystemQuestion) (userQuestion *UserSystemQuestion, e error) {
	if len(LsSystemQuestion) == 0 {
		return nil, ErrorLog("The input systemquestion is zero")
	}
	//1. 把问题记录到用户表 2.把用户归类，增加一个维度
	lsSystemQuestionID := []string{}
	for i := 0; i < len(LsSystemQuestion); i++ {
		lsSystemQuestionID = append(lsSystemQuestionID, LsSystemQuestion[i].QID)
		//把用户归结到某一个类别
		answer:=""
		for k:=0;k<len(LsSystemQuestion[i].AnswerList);k++{
			// answer=answer+LsSystemQuestion[i].AnswerList[k]
			answer=LsSystemQuestion[i].AnswerList[k]
			AppendUserCustomerize(LsSystemQuestion[i].QID, LsSystemQuestion[i].Question,answer, UID )
		}
		// AppendUserCustomerize(LsSystemQuestion[i].QID, LsSystemQuestion[i].Question,answer, UID )
	}
	userQuestion, err := AddUserSystemQuestionID(UID, lsSystemQuestionID)
	if err != nil {
		return nil, err
	}
	return userQuestion, nil
}



//用户回答活动系统问题，这个在用户领取红包时候内部调用，把活动问题传过来
func AnswerActivitySQuestion(UID,AID string, LsSystemQuestion []SystemQuestion) ( e error) {
	if len(LsSystemQuestion) == 0 {
		return ErrorLog("The input systemquestion is zero")
	}
	//1. 把问题记录到用户表 2.把用户归类，增加一个维度
	lsSystemQuestionID := []string{}
	for i := 0; i < len(LsSystemQuestion); i++ {
		lsSystemQuestionID = append(lsSystemQuestionID, LsSystemQuestion[i].QID)
		//把用户归结到某一个类别
		answer:=""
		for k:=0;k<len(LsSystemQuestion[i].AnswerList);k++{
			// answer=answer+LsSystemQuestion[i].AnswerList[k]
			answer=LsSystemQuestion[i].AnswerList[k]
			AppendActivityCustomerize(LsSystemQuestion[i].QID, LsSystemQuestion[i].Question,answer, UID,AID,LsSystemQuestion[i].CID )
		}
		
	}
	return  nil
}




//把用户回答的问题记录在用户问题表中，避免反复提问同样的问题
func AddUserSystemQuestionID(UID string, LsSystemQuestionID []string) (userQuestion *UserSystemQuestion, e error) {
	// AnswerCountDex int
	// AnswerTotalCount int
	directWrite := false
	var err error
	if len(LsSystemQuestionID) == 0 {
		return nil, ErrorLog("Add user System Question ID Fail, the length is 0!\n")
	}
	systemQuestionIdDB := &SystemQuestionId{}
	//drag out the system question ids
	err = db.ShareLock(constant.Hash_SystemQuestionID, constant.KEY_SYSQuestion_Normal, systemQuestionIdDB)
	if err != nil {
		return nil, ErrorLog("Pull the System Question ID fail")
	}
	//pull out the ueserSystemQuestion
	user := &UserSystemQuestion{}
	if err = db.WriteLock(constant.Hash_UserSystemQuestion, UID, user); err != nil {
		//Direct write
		user.LsSystemQuestionAnswered = []string{}
		user.LsSystemQuestionWaitting = []string{}
		user.AnswerTotalCount = 0
		user.AnswerCountDex = 0
		directWrite = true
	}
	user.AnswerCountDex = user.AnswerCountDex + len(LsSystemQuestionID)
	user.AnswerTotalCount = user.AnswerTotalCount + len(LsSystemQuestionID)
	for i := 0; i < len(LsSystemQuestionID); i++ {
		user.LsSystemQuestionAnswered = append(user.LsSystemQuestionAnswered, LsSystemQuestionID[i])
	}

	if systemQuestionIdDB.SystemQuestionValidCount < user.AnswerCountDex {
		user.AnswerCountDex = user.AnswerCountDex%systemQuestionIdDB.SystemQuestionValidCount
	}
	//Arrange the new question
	UpdateUserSystemQuestion(user, 5)
	if directWrite {
		err = db.DirectWrite(constant.Hash_UserSystemQuestion, UID, user)
	} else {
		err = db.WriteBack(constant.Hash_UserSystemQuestion, UID, user)
	}
	if err != nil {
		return nil, ErrorLog("The user write back fail with uid=%s\n", UID)
	}
	return user, nil
}
