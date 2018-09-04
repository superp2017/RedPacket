package common

import (
	. "JsLib/JsLogger"
	"strings"
)

/*
获取地址用户的网点ID
*/
func GetUserDistributionIDLonLat(oLat, oLon float64, oDSI []DistributionSI) (bool, string) {
	var minDistance float64
	minDistance = 10000000000
	for _, v := range oDSI {
		dis := EarthDistance(oLat, oLon, v.Lat, v.Lon)
		if minDistance > dis {
			minDistance = dis
		}
		if dis <= v.Range {
			return true, v.DID
		}
	}
	ErrorLog("The distance with minDistance=%d is larger than range", minDistance)
	return false, ""
}

/*
获取区域用户的网点ID
*/
func GetUserDistributionIDAreaBk(AreaName string, oDSI []DistributionSI) (bool, string) {
	for _, v := range oDSI {
		if strings.Contains(v.Address, AreaName) {
			return true, v.DID
		}
	}
	return false, ""
}

/*
获取区域用户的网点ID
*/
func GetUserDistributionIDArea(AreaName string, oDSI []DistributionSI) (bool, string) {
	for _, v := range oDSI {
		for _, area := range v.LsArea {
			if AreaName != "" && area != "" && strings.Contains(AreaName, area) {
				return true, v.DID
			}
		}
	}
	return false, ""
}
