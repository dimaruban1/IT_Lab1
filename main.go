package main

import (
	"myDb/endpoints"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// SysCatalog.Tables = procedures.LoadTables(params.TableDefaultFilename)
	// printTables()
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
	router.GET("/db/:dbName", endpoints.GetDb)
	router.POST("/db/:dbName", endpoints.CreateDb)

	router.GET("/tables", endpoints.GetSimplifiedTables)
	router.GET("/tables/:name", endpoints.GetTable)
	router.GET("/records/:name", endpoints.GetTableRecords)
	router.GET("/records/:name/:pk", endpoints.GetTableRecord)
	router.GET("/records-project/:name", endpoints.GetProjectedTableRecords)

	router.POST("/tables", endpoints.CreateTable)
	router.POST("/records/:name", endpoints.CreateRecord)

	router.PUT("/records/:name/:pk", endpoints.AlterRecord)

	router.DELETE("/tables/:name", endpoints.DeleteTable)
	router.DELETE("/records/:name/:pk", endpoints.DeleteRecord)

	// Start the Gin server on port 8080
	router.Run(":8080")
}
