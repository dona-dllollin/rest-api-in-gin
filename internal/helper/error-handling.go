package helper

import "github.com/gin-gonic/gin"

// JSONError mengirim response error
func JSONError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
	c.Abort() // optional: hentikan chain handler
}
