package util

import (
	"JsLib/JsConfig"
	"JsLib/JsLogger"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/coocood/freecache"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// const (
// 	appKey    = "appkey-alidayu"
// 	secretKey = "secretkey-alidayu"
// 	vUrl      = "verify-mobile-url-aliday"
// 	vJsnUrl   = "verify-mobile-url-aliday-jsnum"
// )

// var g_appkey string
// var g_secretkey string
// var g_url string
// var g_jsnum_url string

var g_smscache *freecache.Cache

var g_rand_chan chan int

var g_log *JsLogger.ST_Logger

func init() {
	var ok bool
	g_log, ok = JsLogger.GetComLogger()
	if !ok {
		log.Fatalln("can not find common logger.")
	}

	// var e error
	// v, e := vgssdb.Get(appKey)
	// if e != nil {
	// 	log.Fatalf("can not find %s\n", appKey)
	// }
	// g_appkey = string(v)

	// v, e = vgssdb.Get(secretKey)
	// if e != nil {
	// 	log.Fatalf("can not find %s\n", secretKey)
	// }
	// g_secretkey = string(v)

	// v, e = vgssdb.Get(vUrl)
	// if e != nil {
	// 	log.Fatalf("can not find %s\n", vUrl)
	// }
	// g_url = string(v)

	// v, e = vgssdb.Get(vJsnUrl)
	// if e != nil {
	// 	log.Fatalf("can not find %s\n", vJsnUrl)
	// }

	// g_jsnum_url = string(v)

	g_smscache = freecache.NewCache(32 * 1024 * 1024) // 32MB

	g_rand_chan = make(chan int)

	go randCoolie()
}

func randCoolie() {
	rand_gen := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	for {
		g_rand_chan <- rand_gen.Int()
	}
}

func getCode() string {
	code := <-g_rand_chan
	ret := fmt.Sprintf("%06d", code%1000000)
	return ret
}

func verify(signName, code, product, mobile, smscode string) {
	//	verify("注册验证", code, product, mobile, "SMS_3030053"
	para := "?appkey="
	para += JsConfig.CFG.MobileVerify.AppKey
	para += "&secretkey="
	para += JsConfig.CFG.MobileVerify.SecretKey
	para += "&signname="
	para += signName
	para += "&code="
	para += code
	para += "&product="
	para += product
	para += "&mobile="
	para += mobile
	para += "&smscode="
	para += smscode
	//fmt.Printf("%s\n", g_url+para)
	response, e := http.Get(JsConfig.CFG.MobileVerify.VUrl + para)
	b := make([]byte, 2048)
	response.Body.Read(b)

	defer response.Body.Close()
	if e != nil {
		b := make([]byte, 2048)
		n, _ := response.Body.Read(b)
		g_log.Error("verify %s error, rsp:%s\n", mobile, string(b[:n]))
	}
}

func RegisterAuth(mobile, product string, expire int) {
	code := getCode()
	g_smscache.Set([]byte(mobile), []byte(code), expire)
	g_log.Info("---------------------------------------------------------Code=%s\n", code)
	verify("注册验证", code, product, mobile, "SMS_3030053")
}

func IDAuth(mobile string, product string, expire int) {
	code := getCode()
	g_smscache.Set([]byte(mobile), []byte(code), expire)
	verify("君赛科技", code, product, mobile, "SMS_86945210")
}

func LoginAuth(mobile string, product string, expire int) {
	code := getCode()
	g_smscache.Set([]byte(mobile), []byte(code), expire)

	verify("登录验证", code, product, mobile, "SMS_3030055")
}

func LoginExceptionAuth(mobile string, product string, expire int) {
	code := getCode()
	g_smscache.Set([]byte(mobile), []byte(code), expire)

	verify("登录验证", code, product, mobile, "SMS_3030054")
}

// func ActivityAuth(mobile string, product string, expire int) {
// 	code := getCode()
// 	g_smscache.Set([]byte(mobile), []byte(code), expire)

// 	verify("活动验证", code, product, mobile, "SMS_3030052")
// }

func ChangePwdAuth(mobile string, product string, expire int) {
	code := getCode()
	g_smscache.Set([]byte(mobile), []byte(code), expire)
	verify("变更验证", code, product, mobile, "SMS_3030051")
}

func DataChangeAuth(mobile string, product string, expire int) {
	code := getCode()
	g_smscache.Set([]byte(mobile), []byte(code), expire)
	verify("变更验证", code, product, mobile, "SMS_3030050")
}

func JsNumberNotify(field, jsn, pwd, mobile string) {

	para := "?appkey="
	para += JsConfig.CFG.MobileVerify.AppKey
	para += "&secretkey="
	para += JsConfig.CFG.MobileVerify.SecretKey
	para += "&signname="
	para += "君赛认证"
	para += "&mobile="
	para += mobile
	para += "&smscode="
	para += "SMS_86945210"
	para += "&field="
	para += field
	para += "&product="
	para += "君赛科技"
	para += "&name="
	para += jsn
	para += "&pwd="
	para += pwd
	//fmt.Printf("%s\n", g_url+para)

	response, _ := http.Get(JsConfig.CFG.MobileVerify.VJsnUrl + para)
	b := make([]byte, 2048)
	n, e := response.Body.Read(b)

	if e != nil {
		log.Fatalln(e.Error())
	}
	log.Panicf("b = %s\n", string(b[0:n]))

	defer response.Body.Close()
}

func VerifySmsCode(mobile, code string) bool {
	vCode, e := g_smscache.Get([]byte(mobile))
	if e == nil && string(vCode) == code {
		return true
	} else {
		return false
	}
}

////////////////////////////////////////////////////////////////////////////////
//
//新接口
//                                                                            //
////////////////////////////////////////////////////////////////////////////////

const v5_url = "http://www.api.zthysms.com/sendSms.do"
const v5_username = "shxyhy"
const v5_password = "9BApAi"

func verify_ex(code, product, mobile string) {
	tkey := time.Now().Format("20060102150405")

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(v5_password))
	cipherStr := md5Ctx.Sum(nil)

	md5Ctx = md5.New()
	md5Ctx.Write([]byte(hex.EncodeToString(cipherStr) + tkey))

	pwd := hex.EncodeToString(md5Ctx.Sum(nil))
	para := "?username=" + v5_username
	para += "&tkey=" + tkey
	para += "&password=" + pwd
	para += "&mobile=" + mobile
	para += "&content=hello"

	response, e := http.Get(v5_url + para)
	b := make([]byte, 2048)
	response.Body.Read(b)

	defer response.Body.Close()
	if e != nil {
		b := make([]byte, 2048)
		n, _ := response.Body.Read(b)
		g_log.Error("verify %s error, rsp:%s\n", mobile, string(b[:n]))
	}
}

func RegisterAuth_ex(mobile, product string, expire int) {
	code := getCode()
	g_smscache.Set([]byte(mobile), []byte(code), expire)
	g_log.Info("---------------------------------------------------------Code=%s\n", code)
	verify_ex(code, product, mobile)
}
