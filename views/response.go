// views/response.go

package views

import "github.com/gin-gonic/gin"

func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}
