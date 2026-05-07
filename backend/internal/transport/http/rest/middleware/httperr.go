package middleware

import (
	"github.com/Jaeun-Choi98/container-bay/internal/logger"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/http-utils/httperr"
	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, ctxErr := range c.Errors {
			if err, ok := ctxErr.Err.(*httperr.HttpError); ok {
				// if err.Code == http.StatusInternalServerError {
				// 	logger.Infoln(err.ErrMsg)
				// }

				if err.ErrMsg != "" {
					logger.Infof("[REST] error message: %s", err.ErrMsg)
				}
				if err.BaseResponse != nil {
					logger.Infof("[REST] result code: %d", err.BaseResponse.Result)
				}
				c.JSON(err.Code, err.BaseResponse)
			}
		}
	}
}
