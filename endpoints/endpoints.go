package endpoints

import (
	"fmt"
	"myDb/params"
	"myDb/parser"
	"myDb/procedures"
	SysCatalog "myDb/system_catalog"
	"myDb/types"
	"myDb/utility"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetDb(c *gin.Context) {
	dbName := c.Param("dbName")
	SysCatalog.Tables = procedures.LoadTables(params.SaveDir + "\\" + dbName + ".bin")
	if SysCatalog.Tables == nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	params.TableDefaultFilename = params.SaveDir + "\\" + dbName + ".bin"
	c.JSON(http.StatusOK, getSimplifiedTables())
}

func CreateDb(c *gin.Context) {
	dbName := c.Param("dbName")
	utility.CreateFileIfNotExists(params.SaveDir + "\\" + dbName + ".bin")
	SysCatalog.Tables = procedures.LoadTables(params.SaveDir + "\\" + dbName + ".bin")

	params.TableDefaultFilename = params.SaveDir + "\\" + dbName + ".bin"
	c.JSON(http.StatusOK, gin.H{})
}

func GetSimplifiedTables(c *gin.Context) {
	userTables := getSimplifiedTables()
	c.JSON(http.StatusOK, userTables)
}

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
	err := createTable(tableReceived, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, nil)
}

// Delete a table by name
func DeleteTable(c *gin.Context) {
	name := c.Param("name")

	table := SysCatalog.GetTableByName(name)
	filename := params.SaveDir + "\\" + table.DataFileName
	os.Remove(filename)

	if SysCatalog.GetTableByName(name) == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		return
	}

	err := SysCatalog.DeleteTableByName(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	procedures.SaveAllTablesBin(SysCatalog.Tables, params.TableDefaultFilename)

	c.Status(http.StatusNoContent)
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
	fieldMaps := getRecords(table, filename)

	for _, field := range table.Fields {
		if field.Type == types.ColorInvl_t {
			for _, record := range fieldMaps {
				for id, value := range record {
					if id == field.FieldId {
						fmt.Print(value)
						record[id] = ColorInvlToString(value.(map[string]interface{}))
					}
				}
			}
		}
	}
	c.JSON(http.StatusOK, fieldMaps)
}

func ColorInvlToString(c map[string]interface{}) string {
	return fmt.Sprintf("%s:%s:%f", c["color1"], c["color2"], c["interval"])
}

func GetTableRecord(c *gin.Context) {
	name := c.Param("name")
	pk := c.Param("pk")
	fmt.Print(pk)

	table := SysCatalog.GetTableByName(name)
	var pkValue interface{} = pk
	pkField := table.GetPkField()

	if pkField.Type != types.Char_t && pkField.Type != types.String_t {
		pkValue, _ = parser.ParseFieldValue(pkField, pk)
	}
	filename := params.SaveDir + "\\" + table.DataFileName
	fmt.Printf("%s", filename)
	file, _ := os.Open(filename)

	fieldMap := getRecord(table, file, pkValue)
	if fieldMap == nil {
		c.JSON(http.StatusBadRequest, nil)
	}
	// for _, field := range table.Fields {
	// 	if field.Type == types.ColorInvl_t {
	// 		for _, f := range fieldMap {
	// 			if f.ID == field.FieldId {
	// 				f.Value = ColorInvlToString(f.Value.(map[string]interface{}))
	// 			}
	// 		}
	// 	}
	// }
	c.JSON(http.StatusOK, fieldMap)
}

func GetProjectedTableRecords(c *gin.Context) {
	name := c.Param("name")
	fieldNames := make([]string, 0)

	if err := c.ShouldBindJSON(&fieldNames); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	table := SysCatalog.GetTableByName(name)
	file, err := os.Open(params.SaveDir + "\\" + table.DataFileName)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
	}

	fields := make([]types.Field, 0)
	for _, fieldName := range fieldNames {
		f := table.GetFieldByName(fieldName)
		if f == nil {
			continue
		}
		fields = append(fields, *f)
	}
	records := getProjectedTableRecords(file, fields)

	c.JSON(http.StatusOK, records)
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
	fmt.Print(fieldValues)
	record := createRecord(table, fieldValues, filename)
	if record == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no create table >("})
		return
	}

	c.JSON(http.StatusOK, record)
}

func AlterRecord(c *gin.Context) {
	tableName := c.Param("name")
	pk := c.Param("pk")
	table := SysCatalog.GetTableByName(tableName)
	filename := params.SaveDir + "\\" + table.DataFileName

	var pkValue interface{} = pk
	pkField := table.GetPkField()
	if pkField.Type != types.Char_t && pkField.Type != types.String_t {
		pkValue, _ = parser.ParseFieldValue(pkField, pk)
	}

	var newfieldValues map[int32]interface{}
	if err := c.ShouldBindJSON(&newfieldValues); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(newfieldValues) != len(table.Fields) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect field value count"})
		return
	}

	err := alterRecord(table, filename, pkValue, newfieldValues)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, nil)
}

func DeleteRecord(c *gin.Context) {
	tableName := c.Param("name")
	pk := c.Param("pk")
	table := SysCatalog.GetTableByName(tableName)

	var pkValue interface{} = pk
	pkField := table.GetPkField()
	if pkField.Type != types.Char_t && pkField.Type != types.String_t {
		pkValue, _ = parser.ParseFieldValue(pkField, pk)
	}

	filename := params.SaveDir + "\\" + table.DataFileName
	fmt.Println()
	fmt.Print("'")
	fmt.Print(pkValue)
	fmt.Print("'")
	fmt.Println()
	err := deleteRecord(table, filename, pkValue)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, nil)
}
