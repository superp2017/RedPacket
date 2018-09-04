package common

import (
	"constant"
	"db"
	"time"
)

////启动定时任务
func StartTimer() {
	go timeTask(forreachOrder, 2)
}

////定时任务
func timeTask(task func(), hour int) {
	for {
		//////////定时任务//////////////////
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), hour, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		task()
	}
}

///遍历所有的订单，进行无效判断
func forreachOrder() {
	list := getglobalorderlist()
	for _, v := range list {
		data := &ST_Order{}
		if err := db.WriteLock(constant.Hash_Order, v, data); err != nil {
			continue
		}
		if data.Status == constant.C_PAY_WAITPAY {
			if time.Now().Unix()-data.SubmitStamp > 24*3600 {
				data.Status = constant.C_PAY_INVALID
				go Append2Invalid(v)
				go removeFromGlobalOrder(v)
				go removeFormUserOrder(v, data.UID)
			}
		}
	}
}
