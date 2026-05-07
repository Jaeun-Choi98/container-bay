package middleware

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func SpaHandlerRoot(staticPath, indexPath string) gin.HandlerFunc {

	return func(c *gin.Context) {
		url := c.Request.URL.Path
		ext := filepath.Ext(url)

		// 1. 민감한 확장자 차단
		blockedExt := map[string]bool{
			".git": true,
			".ini": true,
			".svg": true,
			".txt": true,
		}
		if blockedExt[ext] {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// 2. 경로 정리 및 탐색 방지
		cleanPath := filepath.Clean(url)

		// 리다이렉트 경로 추가
		redirectPath := map[string]bool{
			filepath.FromSlash("/dashboard"):        true,
			filepath.FromSlash("/volume-directory"): true,
			filepath.FromSlash("/page3"):            true,
			filepath.FromSlash("/page4"):            true,
		}

		// /로 리다이렉트
		if redirectPath[cleanPath] {
			c.Redirect(http.StatusFound, "/")
			return
		}

		absPath := filepath.Join(staticPath, cleanPath)
		_, err := os.Stat(absPath)
		/**
		 * err != nil 일 때,
		 * 1. 파일이 존재x  -> fallback
		 * 2. 파일에 접근 권한이 없음
		 * 3. 드라이브가 연결x
		 * 4. 깨진 심볼릭 링크를 참조
		 *
		 * 1번의 경우에만 Fallback
		 */
		if err != nil {
			// 1. 파일이 존재하지x
			if errors.Is(err, fs.ErrNotExist) {

				c.AbortWithStatus(http.StatusNotFound)
				return
				// 아래는 fallback
				// http.ServeFile(c.Writer, c.Request, filepath.Join(staticPath, indexPath))
				// return
			} else {
				http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		http.FileServer(http.Dir(staticPath)).ServeHTTP(c.Writer, c.Request)
	}
}

func SpaHandlerOther(urlPrefix, staticPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		ext := filepath.Ext(url)

		blockedExt := map[string]bool{
			".git": true,
			".ini": true,
			".txt": true,
		}
		if blockedExt[ext] {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		cleanPath := filepath.Clean(url)
		cleanPath, hasPrefix := strings.CutPrefix(cleanPath, filepath.FromSlash(urlPrefix))
		if !hasPrefix {
			c.Next()
			return
		}

		absPath := filepath.Join(staticPath, cleanPath)

		_, err := os.Stat(absPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			} else {
				http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
				c.Abort()
				return
			}
		}

		http.StripPrefix(urlPrefix, http.FileServer(http.Dir(staticPath))).ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
