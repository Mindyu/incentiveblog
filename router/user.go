package router

import (
	"incentiveblog/utils"
	"net/http"

	"github.com/labstack/echo"
)

//统一响应消息
func ResponseData(c echo.Context, resp *utils.Resp) {
	resp.ErrMsg = utils.GetRecodeText(resp.Errno)
	c.JSON(http.StatusOK, resp)
}

func Pong(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = "0"
	defer ResponseData(c, &resp)
	return nil
}
