package ider

import (
	"JsLib/JsUuid"
)

func GenID() string {
	id := ""
	id = JsUuid.NewV4().String()
	for {
		if id == "" {
			id = JsUuid.NewV4().String()
		} else {
			break
		}
	}
	return id
}

func GenOrderId() string {
	id := ""
	id = JsUuid.NewV4().String()
	for {
		if id == "" {
			id = JsUuid.NewV4().String()
		} else {
			break
		}
	}
	return id[:32]
}
