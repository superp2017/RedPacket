package db

import (
	"JsLib/JsDBCache/ssdb"
	"constant"
	"log"
)

func Unsafe_Keys(db, table string) map[string]string {
	log.Printf("db = %s, table = %s", db, table)
	keys, e := ssdb.MultiHgetAll(db, table)
	ret := make(map[string]string)
	if e == nil {
		for k, _ := range keys {
			ret[k] = "1"
		}
	}
	return ret
}

func GetKeys(table string) []string {
	keys, e := ssdb.MultiHgetAll(constant.C_DB, table)
	ret := []string{}
	if e == nil {
		for k, _ := range keys {
			ret = append(ret, k)
		}
	}
	return ret
}
