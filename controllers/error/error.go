package error

import "bee-structure/controllers/base"

type ErrorController struct {
	base.BaseController
}

//404错误处理
func (c *ErrorController) Error404() {
	msg := c.ReturnMsg("error_404")
	c.ResJson(base.NotFoundCode, msg, nil)
}

//数据库
func (c *ErrorController) ErrorDb() {
	msg := c.ReturnMsg("error_database")
	c.ResJson(base.NotFoundCode, msg, nil)
}
