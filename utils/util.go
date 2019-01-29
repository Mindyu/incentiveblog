package utils

type Resp struct {
	Errno  string `json:"errno"`
	ErrMsg string `json:"errmsg"`
	Data   string `json:"data"`
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
