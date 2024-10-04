package main

import (
	"myDb/endpoints"
	"myDb/params"
	"myDb/procedures"
	SysCatalog "myDb/system_catalog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	SysCatalog.Tables = procedures.LoadTables(params.TableDefaultFilename)
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

	router.GET("/tables", endpoints.GetTables)
	router.GET("/tables/:name", endpoints.GetTable)
	router.GET("/records/:name", endpoints.GetTableRecords)
	router.GET("/records-project", endpoints.GetProjectedTableRecords)

	router.POST("/create-table", endpoints.CreateTable)
	router.POST("/create-record/:name", endpoints.CreateRecord)

	router.PUT("/tables/:id", endpoints.UpdateTable)
	router.PUT("/records/:id", endpoints.UpdateTable)

	router.DELETE("/tables/:name", endpoints.DeleteTable)
	router.DELETE("/records/:name/:pk", endpoints.DeleteRecord)

	// Start the Gin server on port 8080
	router.Run(":8080")

}
