package endpoints

import (
	"fmt"
	"myDb/params"
	"myDb/procedures"
	recording "myDb/records"
	SysCatalog "myDb/system_catalog"
	"myDb/types"
	"myDb/utility"
	"os"
)

func getSimplifiedTables() []types.SimplifiedTable {
	userTables := make([]types.SimplifiedTable, 0)
	for _, table := range SysCatalog.Tables {
		userTables = append(userTables, *types.CastToUserTable(table))
	}
	return userTables
}

// "" means default filename
func createTable(simplifiedTable types.SimplifiedTable, filename string) error {
	if SysCatalog.GetTableByName(simplifiedTable.Name) != nil {
		return fmt.Errorf("table with this name already exists")
	}

	id := len(SysCatalog.Tables)
	newTable := types.CastFromSimplifiedTable(id, simplifiedTable)
	utility.CreateFileIfNotExists(params.SaveDir + newTable.DataFileName)
	SysCatalog.Tables = append(SysCatalog.Tables, *newTable)
	if filename == "" {
		procedures.SaveAllTablesBin(SysCatalog.Tables, params.TableDefaultFilename)
	} else {
		procedures.SaveAllTablesBin(SysCatalog.Tables, filename)
	}

	return nil
}

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

func getRecords(table *types.Table, filename string) []map[int32]interface{} {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}

	records := recording.GetTableRecords(file)
	formattedRecords := getFieldMap(records)
	return formattedRecords
}

func validateTable() {

}

func validateRecord() {

}
