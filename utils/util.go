package utils

import (
	"net/http"

	"github.com/labstack/echo"
)

const (
	RECODE_OK         = "0"
	RECODE_DBERR      = "4001"
	RECODE_SESSIONERR = "4002"
	RECODE_LOGINERR   = "4003"
	RECODE_PARAMERR   = "4004"
	RECODE_UNKNOWERR  = "4101"
)

var recodeText = map[string]string{
	RECODE_OK:         "成功",
	RECODE_DBERR:      "数据库操作错误",
	RECODE_SESSIONERR: "用户未登录",
	RECODE_LOGINERR:   "用户登录失败",
	RECODE_PARAMERR:   "参数错误",
	RECODE_UNKNOWERR:  "未知错误",
}

func RecodeText(code string) string {
	str, ok := recodeText[code]
	if ok {
		return str
	}
	return recodeText[RECODE_UNKNOWERR]
}

type Resp struct {
	Errno  string      `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

//resp数据响应
func ResponseData(c echo.Context, resp *Resp) {
	resp.ErrMsg = RecodeText(resp.Errno)
	c.JSON(http.StatusOK, resp)
}
