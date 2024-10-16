package recording

import (
	"encoding/json"
	"fmt"
	"io"
	"myDb/types"
	"os"
)

func GetTableRecords(file *os.File) [][]types.UserFieldValue {
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

	return allRecords
}

func InsertTableRecord(file *os.File, objectFieldValues []types.UserFieldValue) error {
	var allRecords [][]types.UserFieldValue

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if len(fileBytes) > 0 {
		// Decode existing records if file is not empty
		err = json.Unmarshal(fileBytes, &allRecords)
		if err != nil {
			return err
		}
	}

	// Append the new record (a []FieldValue)
	allRecords = append(allRecords, objectFieldValues)

	// Write the updated data back to the file
	file.Seek(0, 0) // Rewind to the beginning of the file
	err = json.NewEncoder(file).Encode(allRecords)
	if err != nil {
		return err
	}
	return nil
}

func AlterRelationRecord(filename string, pkId int32, pkValue interface{}, newFieldValues []types.UserFieldValue) error {
	err := DeleteTableRecord(filename, pkId, pkValue)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	err = InsertTableRecord(file, newFieldValues)
	if err != nil {
		return err
	}
	file.Close()

	return nil
}

func DeleteTableRecord(filename string, pkId int32, pkValue interface{}) error {
	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	allRecords := GetTableRecords(file)
	file.Close()

	var i1 int32 = -1
	for i, record := range allRecords {
		for _, fieldValue := range record {
			if fieldValue.ID != pkId {
				continue
			}
			var t interface{} = pkValue

			if value, ok := pkValue.(int); ok {
				t = float64(value)
			}
			if value, ok := pkValue.(int32); ok {
				t = float64(value)
			}

			if fieldValue.ID == pkId && t == fieldValue.Value {
				i1 = int32(i)
			}
		}
	}
	if i1 == -1 {
		return fmt.Errorf(`field value not found`)
	}
	allRecords = append(allRecords[:i1], allRecords[i1+1:]...)

	file, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	file.Seek(0, 0)
	err = json.NewEncoder(file).Encode(allRecords)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}
