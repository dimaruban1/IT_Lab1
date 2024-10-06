package endpoints

import (
	"myDb/params"
	"myDb/procedures"
	recording "myDb/records"
	SysCatalog "myDb/system_catalog"
	"myDb/types"
	"myDb/utility"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Get all tables
func GetSimplifiedTables(c *gin.Context) {
	userTables := getSimplifiedTables()
	c.JSON(http.StatusOK, userTables)
}

// Get a specific table by name
func GetTable(c *gin.Context) {
	name := c.Param("name")

	table := SysCatalog.GetTableByName(name)
	if table == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table not found"})
		return
	}

	c.JSON(http.StatusOK, *types.CastToUserTable(*table))
}

// Create a new table
func CreateTable(c *gin.Context) {
	var tableReceived types.SimplifiedTable

	if err := c.ShouldBindJSON(&tableReceived); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createTable(tableReceived, "")

	c.JSON(http.StatusCreated, nil)
}

// Update an existing table
func UpdateTable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedTable types.Table
	if err := c.ShouldBindJSON(&updatedTable); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingTable := SysCatalog.GetTableById(updatedTable.Id)
	if existingTable == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		return
	}

	updatedTable.Id = int32(id)
	SysCatalog.DeleteTableByName(existingTable.Name)
	SysCatalog.Tables = append(SysCatalog.Tables, updatedTable)

	procedures.SaveAllTablesBin(SysCatalog.Tables, params.TableDefaultFilename)

	c.JSON(http.StatusOK, updatedTable)
}

// Delete a table by name
func DeleteTable(c *gin.Context) {
	name := c.Param("name")

	if SysCatalog.GetTableByName(name) == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		return
	}

	SysCatalog.DeleteTableByName(name)
	procedures.SaveAllTablesBin(SysCatalog.Tables, params.TableDefaultFilename)

	c.Status(http.StatusNoContent)
}

func castRecords(records []map[types.Field]*types.FieldValue) [][]types.UserFieldValue {
	userRecords := make([][]types.UserFieldValue, 0)
	for i, record := range records {
		userRecords = append(userRecords, make([]types.UserFieldValue, 0))
		for _, fv := range record {
			userRecords[i] = append(userRecords[i], *types.CastToUserFieldValue(*fv))
		}
	}
	return userRecords
}

// autism
func castRecords2(records [][]*types.FieldValue) [][]types.UserFieldValue {
	userRecords := make([][]types.UserFieldValue, 0)
	for i, record := range records {
		userRecords = append(userRecords, make([]types.UserFieldValue, 0))
		for _, fv := range record {
			userRecords[i] = append(userRecords[i], *types.CastToUserFieldValue(*fv))
		}
	}
	return userRecords
}

func getFieldMap(records [][]types.UserFieldValue) []map[int32]interface{} {
	good_info := make([]map[int32]interface{}, 0)

	for _, record := range records {
		formattedRecord := make(map[int32]interface{})
		for _, fieldValue := range record {
			formattedRecord[fieldValue.ID] = fieldValue.Value
		}
		good_info = append(good_info, formattedRecord)
	}
	return good_info
}

func GetTableRecords(c *gin.Context) {
	name := c.Param("name")
	table := SysCatalog.GetTableByName(name)
	filename := params.SaveDir + "\\" + table.DataFileName
	fieldMap := getRecords(table, filename)
	c.JSON(http.StatusOK, fieldMap)
}

// func GetProjectedTableRecords(c *gin.Context) {
// 	name := c.Param("table")
// 	fieldNames := c.QueryArray("field")

// 	records := recording.ProjectRecords(name, fieldNames)

// 	c.JSON(http.StatusOK, castRecords2(records))
// }

func CreateRecord(c *gin.Context) {
	tableName := c.Param("name")
	table := SysCatalog.GetTableByName(tableName)

	if table == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table not found"})
		return
	}

	filename := params.SaveDir + "\\" + table.DataFileName
	utility.CreateFileIfNotExists(filename)

	var fieldValues map[int32]interface{}
	if err := c.ShouldBindJSON(&fieldValues); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(fieldValues) != len(table.Fields) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect field value count"})
		return
	}

	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record := createRecord(table, fieldValues, file)

	c.JSON(http.StatusOK, record)
}

func DeleteRecord(c *gin.Context) {
	tableName := c.Param("name")
	pkValue := c.Param("pk")

	table := SysCatalog.GetTableByName(tableName)
	filename := params.SaveDir + "\\" + table.DataFileName
	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, f := range table.Fields {
		if f.Key == types.PrimaryKey {
			recording.DeleteTableRecord(file, int32(i), pkValue)
		}
	}
	c.JSON(http.StatusOK, nil)
}
