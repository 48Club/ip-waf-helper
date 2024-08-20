package server

import (
	"net/http"
	"os"

	"github.com/48Club/ip-waf-helper/database"
	"github.com/48Club/ip-waf-helper/types"
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
	var (
		query = []types.IPWaf{}
		resp  = types.AllIPs{}
	)

	tx := database.Server.Find(&query)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}
	for _, v := range query {
		resp = append(resp, v.IP)
	}
	c.JSON(http.StatusOK, resp)
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
