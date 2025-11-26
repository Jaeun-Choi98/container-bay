package middleware

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/http-utils/httperr"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/http-utils/jwt"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/response"
	"github.com/gin-gonic/gin"
)

func StoreMemberIdToContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accToken := ctx.GetHeader("Authorization")
		if accToken != "" {
			accToken, _ = strings.CutPrefix(accToken, "Bearer ")
			if claims, err := jwt.VaildJwtHS256(accToken); err == nil {
				ctx.Set("memberId", claims.Id)
			}
		}
		// var token string
		// if cookie, err := ctx.Request.Cookie("refToken"); err == nil {
		// 	if token = cookie.Value; token != "" {
		// 		token, _ = strings.CutPrefix(token, "Bearer ")
		// 		if claims, err := jwt.VaildJwtHS256(token); err == nil {
		// 			ctx.Set("memberId", claims.Id)
		// 		}
		// 	}
		// }
		ctx.Next()
	}
}

func CheckSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !slices.Contains([]string{"POST", "PUT", "PATCH"}, ctx.Request.Method) {
			ctx.Next()
			return
		}
		accToken := ctx.GetHeader("Authorization")
		if accToken != "" {
			accToken, _ = strings.CutPrefix(accToken, "Bearer ")
			if _, err := jwt.VaildJwtHS256(accToken); err == nil {
				ctx.Next()
				return
			}
		}
		// var token string
		// if cookie, err := ctx.Request.Cookie("refToken"); err == nil {
		// 	log.Println(cookie.Value)
		// 	if token = cookie.Value; token != "" {
		// 		token, _ = strings.CutPrefix(token, "Bearer ")
		// 		if claims, err := jwt.VaildJwtHS256(token); err == nil {
		// 			if service.CheckRefToken(claims.Id, token) {
		// 				ctx.Next()
		// 				return
		// 			}
		// 		}
		// 	}
		// }
		ctx.Error(httperr.UNAUTHORIZED.Add(fmt.Errorf("[REST] CheckSession, invalid token"), response.INVALID_TOKEN))
		ctx.Abort()
	}
}

func VaildJwtToken(tokenString string) error {
	jwtStr, _ := strings.CutPrefix(tokenString, "Bearer ")
	if _, err := jwt.VaildJwtHS256(jwtStr); err != nil {
		return err
	}
	return nil
}
