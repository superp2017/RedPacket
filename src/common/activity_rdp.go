package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"db"
	"math/rand"
	"time"
	. "util"
)

type ActivityRDP struct {
	UID          string  //领用者id
	Name         string  //客户姓名
	Mobile       string  //手机，联系方式
	HeadImageURL string  //头像路径
	Date         string  //日期
	Money        int     //红包金额
	TimeStap     int64   //时间戳
	Lat          float64 //经度
	Lon          float64 //纬度
	City         string  //城市
	LTime        int     //用户从进入到领红包的时间
	Gowil        int     //购买意愿
	IsCorrect    bool    //是否回答正确
	InRange      bool    //是否在范围内
	IsCoupon     bool    //是否领券
}

type Shared struct {
	ShareID    string //分享者的uid
	TotalMoney int    //抢到的钱
	ShareMoney int    //分享佣金
	Date       string //日期
}

type RedPacketGetRecord struct {
	ActivityID       string                   //活动ID
	RDP              []ActivityRDP            //所有红包列表
	AlgoType         int                      //算法类型
	Probability      int                      //红包概率
	TotalMoney       int                      //总金额
	IsAuto           int                      //是否是自动模式
	IsRealTime       int                      //是否是实时到账
	MoneyArrangeWay  string                   //分配方式
	TotalCustomer    int                      //分配人数
	Remaining        int                      //剩余金额
	SharedMoney      int                      //分享的佣金
	SharedRDP        []Shared                 //分享佣金记录
	FinishedPercent  int                      //完成多少了e.g "62.34%"
	ArrangedCustomer int                      //已经抢红包人数
	RobNum           int                      //抢红包的人，不管有没有抢到钱
	LastNum          int                      //剩余数量
	Scale            []string                 //刻度
	Data             map[string][]ActivityRDP ///排序后的数据
	Date             string                   //创建日期
}

func (ac *RedPacketGetRecord) Algo_base() int {
	if ac.Remaining <= 0 {
		return 0
	}

	r := 0
	if ac.ArrangedCustomer < ac.TotalCustomer {
		people := ac.TotalCustomer - ac.ArrangedCustomer //剩余人数
		Ave := int(ac.Remaining / people)                //相对平均值
		r = randInt(1, Ave)
	} else {
		Ave := int(ac.TotalMoney / ac.TotalCustomer) //总平局值
		r = randInt(1, Ave)
	}

	if r < 0 {
		r = 0
	}

	if ac.Remaining <= r {
		return ac.Remaining
	}
	return r
}

func (ac *RedPacketGetRecord) Algo_lessThan(max int) int {
	r := ac.Algo_base()
	for {
		if r >= max {
			r -= max
		} else {
			break
		}
	}
	if ac.Remaining <= r {
		return ac.Remaining
	}
	return r
}

func (ac *RedPacketGetRecord) Algo_lessThan_120() int {
	return ac.Algo_lessThan(120)
}

func (ac *RedPacketGetRecord) Algo_30() int {
	r := 0
	if ac.ArrangedCustomer <= 30 {
		r = randInt(100, 120)
	} else {
		r = ac.Algo_lessThan(120)
	}
	if ac.Remaining <= r {
		r = ac.Remaining
	}
	return r
}

func (ac *RedPacketGetRecord) Algo_100() int {
	r := 0
	if ac.ArrangedCustomer <= 100 {
		r = randInt(100, 120)
	} else {
		r = ac.Algo_lessThan(120)
	}
	if ac.Remaining <= r {
		r = ac.Remaining
	}
	return r
}

func (ac *RedPacketGetRecord) Algo_Step() int {
	r := 0
	if time.Now().Unix()%2 == 0 {
		r = randInt(100, 120)
	} else {
		r = 0
	}
	if ac.Remaining <= r {
		r = ac.Remaining
	}
	return r
}

func (ac *RedPacketGetRecord) Algo_100_Step() int {
	r := 0
	if ac.ArrangedCustomer < 100 {
		r = randInt(100, 120)
	} else {
		r = ac.Algo_Step()
	}
	if ac.Remaining <= r {
		r = ac.Remaining
	}
	return r
}

///前一半必领，后一半随机
func (ac *RedPacketGetRecord) Algo_half_get() int {
	r := 0
	if ac.ArrangedCustomer <= ac.TotalCustomer {
		r = randInt(100, 120)
	} else {
		r = ac.Algo_lessThan(150)
	}
	if ac.Remaining <= r {
		r = ac.Remaining
	}
	return r
}

func (ac *RedPacketGetRecord) Algo_all_get() int {
	r := 0
	r = randInt(100, 120)
	if ac.Remaining <= r {
		r = ac.Remaining
	}
	return r
}

//获取活动领取记录
func getActivityRDP(Acid string) (*RedPacketGetRecord, error) {
	data := &RedPacketGetRecord{}
	err := db.ShareLock(constant.Hash_ActivityRDP, Acid, data)
	return data, err
}

func GetActivityRDP(session *JsNet.StSession) {
	type INFO struct {
		ActivityID string //活动id
		SortType   int    //排序方式
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ActivityID == "" {
		ForwardEx(session, "1", nil, "GetActivityRDP ActivityID is empty\n")
		return
	}

	re, err := getActivityRDP(st.ActivityID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if re.Data == nil {
		re.Data = make(map[string][]ActivityRDP)
	}

	for _, v := range re.RDP {
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
		if l, ok := re.Data[str]; ok {
			l = append(l, v)
			re.Data[str] = l
		} else {
			L := RDPRecord{v}
			re.Data[str] = L
		}
		AppendUniqueString(&re.Scale, str)
	}

	Forward(session, "0", re)
}

///添加活动的红包领取记录,生成预定义的红包金额
func NewActivityRDP(ac *Activity) error {
	if ac == nil {
		return ErrorLog("NewActivityRDP failed,Activity is nil\n ")
	}
	data := RedPacketGetRecord{}

	err := db.WriteLock(constant.Hash_ActivityRDP, ac.ActivityID, &data)
	if err != nil {
		data = RedPacketGetRecord{
			ActivityID:       ac.ActivityID,
			Date:             CurTime(),
			TotalMoney:       ac.TotalMoney,
			TotalCustomer:    ac.TotalCustomer,
			MoneyArrangeWay:  ac.MoneyArrangeWay,
			AlgoType:         ac.AlgoType,
			Probability:      ac.Probability,
			IsAuto:           ac.IsAuto,
			IsRealTime:       ac.IsRealTime,
			SharedMoney:      ac.ShardMoney,
			ArrangedCustomer: 0,
			FinishedPercent:  0,
			RobNum:           0,
			Remaining:        ac.TotalMoney,
		}
		return db.DirectWrite(constant.Hash_ActivityRDP, ac.ActivityID, &data)
	}

	data.TotalMoney = ac.TotalMoney
	data.TotalCustomer = ac.TotalCustomer
	data.MoneyArrangeWay = ac.MoneyArrangeWay
	data.AlgoType = ac.AlgoType
	data.IsRealTime = ac.IsRealTime
	return db.WriteBack(constant.Hash_ActivityRDP, ac.ActivityID, &data)
}

///生成一个随机整数[min,max]
func randInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max-min) + min
	return randNum
}

///根据概率随机
func randProbability(Probability int) int {
	r := randInt(1, 120)
	d := randInt(1, 100)
	if Probability > d {
		r = 0
	}
	return r
}

////添加红包到活动的领取记录，并检查是否领过

func AppendActivityRDP(user *User, SharedID, ActivityID, Status, City string,
	Lat, Lon float64, flag, DisRemaining, LTime int, isCorrect, inRange bool) (int, int, bool, error) {
	Money := 0
	sharMoney := 0
	IsStop := false
	isShared := (SharedID != "") && (SharedID != user.UID)

	data := RedPacketGetRecord{}
	if err := db.WriteLock(constant.Hash_ActivityRDP, ActivityID, &data); err != nil {
		return sharMoney, Money, IsStop, err
	}
	exist := false
	U_INDEX := -1
	for i, v := range data.RDP {
		if v.UID == user.UID {
			U_INDEX = i
			if v.InRange {
				exist = true
				break
			}
		}
	}

	if exist {
		e := db.WriteBack(constant.Hash_ActivityRDP, ActivityID, &data)
		if e != nil {
			return sharMoney, Money, IsStop, e
		}
		return sharMoney, Money, IsStop, ErrorLog("this User:%s is Recevie  ACtivity:%s RDP\n", user.UID, ActivityID)
	}

	if isCorrect && inRange {
		if Status == constant.Status_Active {
			if data.IsAuto != -1 {
				switch data.AlgoType {
				case 0:
					Money = data.Algo_base()
				case 1:
					Money = data.Algo_lessThan_120()
				case 2:
					Money = data.Algo_30()
				case 3:
					Money = data.Algo_100()
				case 4:
					Money = data.Algo_Step()
				case 5:
					Money = data.Algo_100_Step()
				case 6:
					Money = data.Algo_half_get()
				case 7:
					Money = data.Algo_all_get()
				default:
					Money = data.Algo_base()
				}
			} else {
				Money = randProbability(data.Probability)
			}
		} else {
			Money = 0
		}

		sharMoney = int(Money * 30 / 100)
		if data.IsRealTime == 1 {
			if Money < 100 {
				Money = 0
				sharMoney = 0
			}
		}

		if DisRemaining >= 0 && DisRemaining < Money {
			Money = DisRemaining
			sharMoney = 0
		}

		data.FinishedPercent = data.RobNum * 100 / data.TotalCustomer
		if data.FinishedPercent >= 100 {
			data.FinishedPercent = 100
		}
		///////////////////////////////

		if data.Remaining <= Money {
			Money = data.Remaining
			sharMoney = 0
		}
		if Money > 0 {
			data.ArrangedCustomer++
			data.LastNum = data.TotalCustomer - data.ArrangedCustomer
			data.Remaining -= Money
			if isShared {
				if data.Remaining < sharMoney {
					sharMoney = data.Remaining
				}
				data.Remaining -= sharMoney
				data.SharedMoney += sharMoney
			}
		}
	}

	if data.LastNum <= 0 {
		data.LastNum = 0
	}

	if data.Remaining <= 0 {
		data.Remaining = 0
		if data.ArrangedCustomer >= data.TotalCustomer {
			IsStop = true
		}
	}
	rdp := ActivityRDP{
		UID:          user.UID,
		Money:        Money,
		Name:         user.Nickname,
		Mobile:       user.Mobile,
		HeadImageURL: user.HeadImageURL,
		Date:         CurTime(),
		TimeStap:     CurStamp(),
		Lat:          Lat,
		Lon:          Lon,
		IsCorrect:    isCorrect,
		City:         City,
		InRange:      inRange,
		LTime:        LTime,
		IsCoupon:     false,
	}
	if U_INDEX != -1 {
		data.RDP[U_INDEX] = rdp
	} else {
		data.RDP = append(data.RDP, rdp)
	}
	data.RobNum = len(data.RDP)
	if flag == -1 || !inRange {
		Money = 0
		sharMoney = 0
	}
	if isShared && Money > 0 && sharMoney > 0 {
		data.SharedRDP = append(data.SharedRDP, Shared{
			ShareID:    SharedID,
			TotalMoney: Money,
			ShareMoney: sharMoney,
			Date:       CurTime(),
		})
	}

	return sharMoney, Money, IsStop, db.WriteBack(constant.Hash_ActivityRDP, ActivityID, &data)
}

///更新活动的RDP，券的领用情况
func updateActivityRDP(ActivityID, UID string) error {
	data := RedPacketGetRecord{}
	if err := db.WriteLock(constant.Hash_ActivityRDP, ActivityID, &data); err != nil {
		return err
	}

	for _, v := range data.RDP {
		if v.UID == UID {
			v.IsCoupon = true
		}
	}

	return db.WriteBack(constant.Hash_ActivityRDP, ActivityID, &data)
}
