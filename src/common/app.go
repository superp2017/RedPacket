package common

import (
	"JsLib/JsNet"
	"db"
	. "util"
)

//获取当前版本号
func GetVersion(session *JsNet.StSession) {
	ver, err := getVersion()
	if err != nil {
		ForwardEx(session, "1", ver, err.Error())
		return
	}
	Forward(session, "0", ver)
}

//设置当前版本号
func SetVersion(session *JsNet.StSession) {
	ver := 0
	type INFO struct {
		Version int
	}

	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.Version <= 0 {
		ForwardEx(session, "1", nil, "Version number must than 0\n")
		return
	}

	if ver, err := getVersion(); err == nil {
		if st.Version <= ver {
			ForwardEx(session, "1", nil, "Set version failed,Last version: %d,cur version: %d", ver, st.Version)
			return
		}
	}

	if err := db.Set("Version", st.Version); err != nil {
		ForwardEx(session, "1", ver, err.Error())
		return
	}
	Forward(session, "0", st.Version)
}

func getVersion() (int, error) {
	ver := 0
	err := db.Get("Version", &ver)
	return ver, err
}
