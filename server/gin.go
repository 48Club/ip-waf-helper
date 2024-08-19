package server

import (
	"ip-waf-helper/database"
	"ip-waf-helper/types"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func postHandlerFunc(c *gin.Context) {
	var query = types.IPWaf{}
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create if not exists
	tx := database.Server.FirstOrCreate(&query, query)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, query)
}
func getHandlerFunc(c *gin.Context) {
	var query = []types.IPWaf{}
	tx := database.Server.Find(&query)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, query)
}

func Run() error {
	r := gin.Default()

	token := os.Getenv("GIN_TOKEN")

	r.Use(func(ctx *gin.Context) {
		if ctx.GetHeader("Authorization") != token {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}
	})

	r.POST("/", postHandlerFunc)
	r.GET("/", getHandlerFunc)

	return r.Run(":80")

}
