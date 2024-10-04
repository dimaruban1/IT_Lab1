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
func GetTables(c *gin.Context) {
	userTables := make([]types.UserTable, 0)
	for _, table := range SysCatalog.Tables {
		userTables = append(userTables, *types.CastToUserTable(table))
	}

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

	var tableReceived types.UserTable

	if err := c.ShouldBindJSON(&tableReceived); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if SysCatalog.GetTableByName(tableReceived.Name) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table with this name already exists"})
		return
	}
	id := len(SysCatalog.Tables)
	newTable := types.CastFromUserTable(id, tableReceived)

	utility.CreateFileIfNotExists(params.SaveDir + newTable.DataFileName)
	SysCatalog.Tables = append(SysCatalog.Tables, *newTable)
	procedures.SaveAllTablesBin(SysCatalog.Tables, params.TableDefaultFilename)

	c.JSON(http.StatusCreated, newTable)
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

func GetTableRecords(c *gin.Context) {
	name := c.Param("name")

	records := recording.GetRecords(name)

	c.JSON(http.StatusOK, castRecords(records))
}

func GetProjectedTableRecords(c *gin.Context) {
	name := c.Param("table")
	fieldNames := c.QueryArray("field")

	records := recording.ProjectRecords(name, fieldNames)

	c.JSON(http.StatusOK, castRecords2(records))
}

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

// testing
func createRecord(table *types.Table, fieldValues map[int32]interface{}, file *os.File) []types.UserFieldValue {
	userRecord := make([]types.UserFieldValue, 0)
	for key, value := range fieldValues {
		f := new(types.UserFieldValue)
		f.ID = key
		f.Value = value
		userRecord = append(userRecord, *f)
	}

	recording.InsertTableRecord(file, userRecord)
	return userRecord
}
