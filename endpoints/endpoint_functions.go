package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
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
	// validation of tables
	if err := validateTable(simplifiedTable); err != nil {
		return err
	}

	id := len(SysCatalog.Tables)
	newTable := types.CastFromSimplifiedTable(id, simplifiedTable)
	if newTable == nil {
		return fmt.Errorf("failed to cast table, prob field type")
	}
	utility.CreateFileIfNotExists(params.SaveDir + "\\" + newTable.DataFileName)
	SysCatalog.Tables = append(SysCatalog.Tables, *newTable)
	if filename == "" {
		procedures.SaveAllTablesBin(SysCatalog.Tables, params.TableDefaultFilename)
	} else {
		procedures.SaveAllTablesBin(SysCatalog.Tables, filename)
	}

	return nil
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

func help_me(val interface{}) interface{} {
	var t interface{}
	if value, ok := val.(int); ok {
		t = float64(value)
	} else if value, ok := val.(int32); ok {
		t = float64(value)
	} else if value, ok := val.(map[string]interface{}); ok {
		t = ColorInvlToString(value)
	} else {
		return val
	}
	return t
}

func getRecord(table *types.Table, file *os.File, pkValue interface{}) []types.UserFieldValue {
	file.Seek(0, 0)
	records := recording.GetTableRecords(file)
	fmt.Printf("%d", len(records))
	fmt.Print(pkValue)
	pkFieldId := table.GetPkField().FieldId

	for _, record := range records {
		for _, field := range record {
			t1 := help_me(field.Value)
			t2 := help_me(pkValue)

			if field.ID == pkFieldId && t1 == t2 {
				return record
			}
		}
	}

	return nil
}

func createRecord(table *types.Table, fieldValues map[int32]interface{}, filename string) []types.UserFieldValue {
	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return nil
	}

	pkId := table.GetPkField().FieldId
	for key, value := range fieldValues {
		if key == pkId {
			matchingRecord := getRecord(table, file, value)
			if matchingRecord != nil {
				return nil
			}
		}
	}
	file.Close()

	file, err = os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return nil
	}
	userRecord := make([]types.UserFieldValue, 0)
	for key, value := range fieldValues {
		f := new(types.UserFieldValue)
		f.ID = key
		f.Value = value
		userRecord = append(userRecord, *f)
	}

	err = recording.InsertTableRecord(file, userRecord)
	if err != nil {
		return nil
	}
	file.Close()
	return userRecord
}

func alterRecord(table *types.Table, filename string, pkValue interface{}, newFieldValues map[int32]interface{}) error {
	userRecord := make([]types.UserFieldValue, 0)
	for key, value := range newFieldValues {
		f := new(types.UserFieldValue)
		f.ID = key
		f.Value = value
		userRecord = append(userRecord, *f)
	}
	for _, field := range table.Fields {
		if field.Key == 'P' {
			recording.AlterRelationRecord(filename, field.FieldId, pkValue, userRecord)
		}
	}

	return nil
}

func deleteRecord(table *types.Table, filename string, pkValue interface{}) error {
	err := recording.DeleteTableRecord(filename, table.GetPkField().FieldId, pkValue)
	if err != nil {
		return err
	}
	return nil
}

func getProjectedTableRecords(file *os.File, projectedFields []types.Field) [][]types.UserFieldValue {
	file.Seek(0, 0)
	allRecords := make([][]types.UserFieldValue, 0)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil
	}

	if len(fileBytes) > 0 {
		err = json.Unmarshal(fileBytes, &allRecords)
		if err != nil {
			return nil
		}
	}
	projectedFieldMap := make(map[int32]bool)
	for _, field := range projectedFields {
		projectedFieldMap[field.FieldId] = true
	}

	resultRecords := make([][]types.UserFieldValue, 0)
	for _, record := range allRecords {
		projectedRecord := make([]types.UserFieldValue, 0)
		for _, fieldValue := range record {
			if projectedFieldMap[fieldValue.ID] {
				projectedRecord = append(projectedRecord, fieldValue)
			}
		}

		if len(projectedRecord) > 0 {
			resultRecords = append(resultRecords, projectedRecord)
		}
	}

	return resultRecords
}

func validateTable(table types.SimplifiedTable) error {
	i := 0
	if SysCatalog.GetTableByName(table.Name) != nil {
		return fmt.Errorf("table with this name already exists")
	}
	if len(table.Fields) == 0 {
		return fmt.Errorf("add some fields")
	}
	for _, f := range table.Fields {
		if f.Key == "P" {
			i++
		}
		if f.Key == "P" && f.Type == "colorInvl" {
			return fmt.Errorf("colorInvl can`t be primary key")
		}
		if i > 1 {
			return fmt.Errorf("composite primary keys are not supported")
		}
	}
	if i == 0 {
		return fmt.Errorf("primary key is required")
	}
	return nil
}

func validateRecord() {

}
