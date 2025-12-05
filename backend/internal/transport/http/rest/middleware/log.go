package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/Jaeun-Choi98/container-bay/internal/logger"
	"github.com/gin-gonic/gin"
)

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		// spa로 들어오는 요청은 따로 로그를 남기지 않음.
		if strings.Contains(c.Request.URL.Path, "static") {
			return
		}

		logText := fmt.Sprintf(
			"[%s] %s \"%s %s %s\" %d %.3fsec \"%s\"",
			time.Now().Format(time.RFC1123),
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Proto,
			c.Writer.Status(),
			time.Since(t).Seconds(),
			c.Request.UserAgent(),
		)
		logger.Println(logText)
	}
}
