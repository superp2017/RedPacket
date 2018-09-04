package common

import (
	. "JsLib/JsLogger"
	"constant"
	"db"
	. "util"
	"JsLib/JsNet"
)

//对系统级别，每个答案对应的用户列表
type Customerize struct {
	CustomerizeId    string         //系统问题ID
	SystemQuestionId string         //系统问题索引
	Question         string         //问题
	Answer           string         //答案
	HmCustomer       map[string]int //用户列表
	HashCount        int            //哈希列表的个数，方便后续统计
	TotalCount       int            //总共个数
	AID              string         //活动ID
	CID              string          //公司ID

}

///Add User to the Customerize
func AppendUserCustomerize(SystemQuestionId, Question, Answer, UID string) (e error) {
	type AllCID struct{
		LsAllID []string
	}
	//Judge whethe this is an QuestionAnswer ID
	customerizeID := SystemQuestionId + "&" + Answer

		//save the customerID to the db
	allcid:=&AllCID{}
	err := db.WriteLock(constant.Hash_Customerize, "AllID", allcid)
	if err != nil {
		AppendUniqueString(&allcid.LsAllID,customerizeID)
		db.DirectWrite(constant.Hash_Customerize, "AllID", allcid)
	}else{
		db.WriteBack(constant.Hash_Customerize, "AllID", allcid)
	}

	//pull out the ueser
	customerizeDB := &Customerize{}
	if err := db.WriteLock(constant.Hash_Customerize, customerizeID, customerizeDB); err != nil {
		customerizeDB.CustomerizeId = customerizeID
		customerizeDB.SystemQuestionId = SystemQuestionId
		customerizeDB.Question = Question
		customerizeDB.Answer = Answer
		customerizeDB.HmCustomer = make(map[string]int)
		customerizeDB.HmCustomer[UID] = 1
		customerizeDB.HashCount = 1
		customerizeDB.TotalCount = 1
		if err := db.DirectWrite(constant.Hash_Customerize, customerizeID, customerizeDB); err != nil {
			return ErrorLog("DirectWrite Curimise=%s failed\n", customerizeID)
		}
		return nil
	}

	count, ok := customerizeDB.HmCustomer[UID]
	if !ok {
		customerizeDB.HmCustomer = make(map[string]int)
		customerizeDB.HmCustomer[UID] = 1
		customerizeDB.HashCount = customerizeDB.HashCount + 1

	} else {
		customerizeDB.HmCustomer[UID] = count + 1
	}
	customerizeDB.TotalCount = customerizeDB.TotalCount + 1
	if err := db.WriteBack(constant.Hash_Customerize, customerizeID, customerizeDB); err != nil {
		return ErrorLog("Write back with customerID=%s fail\n", customerizeID)
	}
	return nil
}


///Add Activity to the Customerize
func AppendActivityCustomerize(SystemQuestionId, Question, Answer, UID, AID,CID string) (e error) {
	type AllCID struct{
		LsAllID []string
	}

	type ActivityCID struct{
		HmAID map[string][]string
	}

	
	type CompanyCID struct{
		HmCID map[string][]string
	}
	//Judge whethe this is an QuestionAnswer ID
	customerizeID := AID+SystemQuestionId + "&" + Answer
	//pull out the ueser
	customerizeDB := &Customerize{}
	if err := db.WriteLock(constant.Hash_ActivityCustomerize, customerizeID, customerizeDB); err != nil {
		customerizeDB.CustomerizeId = customerizeID
		customerizeDB.SystemQuestionId = SystemQuestionId
		customerizeDB.Question = Question
		customerizeDB.AID=AID
		customerizeDB.CID=CID
		customerizeDB.Answer = Answer
		customerizeDB.HmCustomer = make(map[string]int)
		customerizeDB.HmCustomer[UID] = 1
		customerizeDB.HashCount = 1
		customerizeDB.TotalCount = 1
		if err := db.DirectWrite(constant.Hash_ActivityCustomerize, customerizeID, customerizeDB); err != nil {
			return ErrorLog("DirectWrite Curimise=%s failed\n", customerizeID)
		}
		// return nil
	}else{
		count, ok := customerizeDB.HmCustomer[UID]
		if !ok {
			customerizeDB.HmCustomer = make(map[string]int)
			customerizeDB.HmCustomer[UID] = 1
			customerizeDB.HashCount = customerizeDB.HashCount + 1
	
		} else {
			customerizeDB.HmCustomer[UID] = count + 1
		}
		customerizeDB.TotalCount = customerizeDB.TotalCount + 1
		if err := db.WriteBack(constant.Hash_ActivityCustomerize, customerizeID, customerizeDB); err != nil {
			return ErrorLog("Write back with customerID=%s fail\n", customerizeID)
		}
	}
	  //save the customerID to the db
		allcid:=&AllCID{}
		err := db.WriteLock(constant.Hash_ActivityCustomerize, "AllID", allcid)
		allcid.LsAllID=AppendUniqueStringEx(allcid.LsAllID,customerizeID)
		if err != nil {
			db.DirectWrite(constant.Hash_ActivityCustomerize, "AllID", allcid)
		}else{
			db.WriteBack(constant.Hash_ActivityCustomerize, "AllID", allcid)
		}
		//save the activityid-customerize id

		acid:=&ActivityCID{}
		err = db.WriteLock(constant.Hash_ActivityCustomerize, "ASID", acid)
		if err != nil {
			acid.HmAID=map[string][]string{}
			acid.HmAID[AID]=[]string{}
			acid.HmAID[AID]=AppendUniqueStringEx(acid.HmAID[AID],customerizeID)
			db.DirectWrite(constant.Hash_ActivityCustomerize, "ASID", acid)
		}else{
			_,ok:=acid.HmAID[AID]
			if !ok{
				acid.HmAID[AID]=[]string{}
			}
			acid.HmAID[AID]=AppendUniqueStringEx(acid.HmAID[AID],customerizeID)
			db.WriteBack(constant.Hash_ActivityCustomerize, "ASID", acid)
		}


				//save the company-customerize id

		ccid:=&CompanyCID{}
			err = db.WriteLock(constant.Hash_ActivityCustomerize, "CSID", ccid)
			if err != nil {
			ccid.HmCID=map[string][]string{}
			ccid.HmCID[CID]=[]string{}
			ccid.HmCID[CID]=AppendUniqueStringEx(ccid.HmCID[CID],customerizeID)
			db.DirectWrite(constant.Hash_ActivityCustomerize, "CSID", ccid)
		  }else{
				_,ok:=ccid.HmCID[CID]
				if !ok{
					ccid.HmCID[CID]=[]string{}
				}
				ccid.HmCID[CID]=AppendUniqueStringEx(ccid.HmCID[CID],customerizeID)
				db.WriteBack(constant.Hash_ActivityCustomerize, "CSID", ccid)
			}
	return nil
}


//获取一个Curimize
func GetCustomerizeInfo(customerizeID string) (*Customerize, error) {
	if customerizeID == "" {
		return nil, ErrorLog("customerize failed,customerizeID is empty !\n")
	}
	customerize := &Customerize{}
	err := db.ShareLock(constant.Hash_Customerize, customerizeID, customerize)
	return customerize, err
}


//获取一个Activity Curimize

func GetACustomerizeInfo(customerizeID string) (*Customerize, error) {
	if customerizeID == "" {
		return nil, ErrorLog("customerize failed,customerizeID is empty !\n")
	}
	customerize := &Customerize{}
	err := db.ShareLock(constant.Hash_ActivityCustomerize, customerizeID, customerize)
	return customerize, err
}

//获取系统所有的用户化列表
func QueryAllCustomrize(session *JsNet.StSession) {
	type AllCID struct{
		LsAllID []string
	}
	allID:=&AllCID{}
	LsCustomrize:=[]*Customerize{}
	//获取所有的ID列表
	err:=db.ShareLock(constant.Hash_Customerize, "AllID", allID)
	if err!=nil{
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	//获取所有的Customize

	for i:=0;i<len(allID.LsAllID);i++{
		customerize:=&Customerize{}
		err:=db.ShareLock(constant.Hash_Customerize, allID.LsAllID[i], customerize)
		if err==nil{
			LsCustomrize=append(LsCustomrize,customerize)
		}
	}
	Forward(session, "0", LsCustomrize)
}

//获取一个公司所有的问题人群列表

func QueryDedicateCompanyCustomrize(session *JsNet.StSession) {
	
		type INFO struct {
			CID string //问题id
		}

		
	type ActivityCID struct{
		HmCID map[string][]string
	}
		st := &INFO{}
		if err := session.GetPara(st); err != nil {
			ForwardEx(session, "1", nil, err.Error())
			return
		}
		//get the cid
		allID:=&ActivityCID{}
		LsCustomrize:=[]*Customerize{}
		//获取所有的ID列表
		err:=db.ShareLock(constant.Hash_ActivityCustomerize, "CSID", allID)
		if err!=nil{
			ForwardEx(session, "1", nil, err.Error())
			return
		}
		_,ok:=allID.HmCID[st.CID]
		if !ok{
			ForwardEx(session, "1", nil, "There is no customrize data for this company")
			return
		}
		//获取所有的Customize
		for i:=0;i<len(allID.HmCID[st.CID]);i++{
			customerize:=&Customerize{}
			err:=db.ShareLock(constant.Hash_ActivityCustomerize, allID.HmCID[st.CID][i], customerize)
			if err==nil{
				LsCustomrize=append(LsCustomrize,customerize)
			}
		}
		Forward(session, "0", LsCustomrize)
	}




//获取一个活动所有的问题人群列表

func QueryDedicateActivityCustomrize(session *JsNet.StSession) {
	
		type INFO struct {
			AID string //问题id
		}

		
	type ActivityCID struct{
		HmAID map[string][]string
	}
	// type CompanyCID struct{
	// 	HmCID map[string][]string
	// }
		st := &INFO{}
		if err := session.GetPara(st); err != nil {
			ForwardEx(session, "1", nil, err.Error())
			return
		}
		//get the aid
		allID:=&ActivityCID{}
		LsCustomrize:=[]*Customerize{}
		//获取所有的ID列表
		err:=db.ShareLock(constant.Hash_ActivityCustomerize, "ASID", allID)
		if err!=nil{
			ForwardEx(session, "1", nil, err.Error())
			return
		}
		_,ok:=allID.HmAID[st.AID]
		if !ok{
			ForwardEx(session, "1", nil, "There is no customrize data for this activity")
			return
		}
		//获取所有的Customize
		for i:=0;i<len(allID.HmAID[st.AID]);i++{
			customerize:=&Customerize{}
			err:=db.ShareLock(constant.Hash_ActivityCustomerize, allID.HmAID[st.AID][i], customerize)
			if err==nil{
				LsCustomrize=append(LsCustomrize,customerize)
			}
		}
		Forward(session, "0", LsCustomrize)
	}


//获取所有公司的问题人群列表

func QueryAllCompanyCustomrize(session *JsNet.StSession) {
	type AllCID struct{
		LsAllID []string
	}
	allID:=&AllCID{}
	LsCustomrize:=[]*Customerize{}
	//获取所有的ID列表
	err:=db.ShareLock(constant.Hash_ActivityCustomerize, "AllID", allID)
	if err!=nil{
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	//获取所有的Customize

	for i:=0;i<len(allID.LsAllID);i++{
		customerize:=&Customerize{}
		err:=db.ShareLock(constant.Hash_ActivityCustomerize, allID.LsAllID[i], customerize)
		if err==nil{
			LsCustomrize=append(LsCustomrize,customerize)
		}
	}
	Forward(session, "0", LsCustomrize)
}


func Restsys(session *JsNet.StSession) {
	type AllCID struct{
		LsAllID []string
	}

	type ActivityCID struct{
		HmAID map[string][]string
	}

	
	type CompanyCID struct{
		HmCID map[string][]string
	}

	allID:=&AllCID{}
	activitycid:=&ActivityCID{}
	companycid:=&CompanyCID{}

	allIDbk:=[]string{}
	activitycidbk:=[]string{}
	companycidbk:=[]string{}


	//获取所有的ID列表
	// -------------------------------
	db.WriteLock(constant.Hash_ActivityCustomerize, "AllID", allID)
	for i:=0;i<len(allID.LsAllID);i++{
		allIDbk=AppendUniqueStringEx(allIDbk,allID.LsAllID[i])
	}
	allID.LsAllID=allIDbk
	db.WriteBack(constant.Hash_ActivityCustomerize, "AllID", allID)
	// -------------------------------
	db.WriteLock(constant.Hash_ActivityCustomerize, "CSID", companycid)
	for k,_:=range companycid.HmCID{
		companycidbk=[]string{}
		for i:=0;i<len(companycid.HmCID[k]);i++{
			companycidbk=AppendUniqueStringEx(companycidbk,companycid.HmCID[k][i])
		}
		companycid.HmCID[k]=companycidbk
	}
	db.WriteBack(constant.Hash_ActivityCustomerize, "CSID", companycid)

		// -------------------------------
		db.WriteLock(constant.Hash_ActivityCustomerize, "ASID", activitycid)
		for k,_:=range activitycid.HmAID{
			activitycidbk=[]string{}
			for i:=0;i<len(activitycid.HmAID[k]);i++{
				activitycidbk=AppendUniqueStringEx(activitycidbk,activitycid.HmAID[k][i])
			}
			activitycid.HmAID[k]=activitycidbk
		}
		db.WriteBack(constant.Hash_ActivityCustomerize, "ASID", activitycid)
		// -------------------------------
	Forward(session, "0", nil)
}
