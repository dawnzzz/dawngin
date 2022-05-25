package dain

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		// 开始时间
		startTime := time.Now()
		// 向后处理请求
		c.Next()
		// 处理结束，输出日志
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(startTime))
	}
}
