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
		absPath := filepath.Join(staticPath, cleanPath)

		_, err := os.Stat(absPath)

		if err != nil {

			if errors.Is(err, fs.ErrNotExist) {
				http.ServeFile(c.Writer, c.Request, filepath.Join(staticPath, indexPath))
				return
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
