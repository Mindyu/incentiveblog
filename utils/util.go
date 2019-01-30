package utils

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

type Resp struct {
	Errno  string      `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

//错误编码定义

const (
	RECODE_OK        = "0"
	RECODE_DBERR     = "1001"
	RECODE_LOGINERR  = "1002"
	RECODE_PARAMERR  = "1003"
	RECODE_UNKONWERR = "1100"
)

var recodeText = map[string]string{
	RECODE_OK:        "成功",
	RECODE_DBERR:     "数据库操作错误",
	RECODE_LOGINERR:  "登陆错误",
	RECODE_PARAMERR:  "请求参数错误",
	RECODE_UNKONWERR: "未知错误",
}

func GetRecodeText(code string) string {
	str, ok := recodeText[code]
	if ok {
		return str
	}
	return recodeText[RECODE_UNKONWERR]
}

//生成token
func MakeToken(orgData1, orgData2 []byte) string {
	timestamp := []byte(fmt.Sprintf("%d", time.Now().Nanosecond()))
	data := bytes.Join([][]byte{orgData1, orgData2, timestamp}, []byte{})
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}
