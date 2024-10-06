package recording

import (
	"encoding/json"
	"io"
	"myDb/types"
	"os"
)

func GetTableRecords(file *os.File) [][]types.UserFieldValue {
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

func GetProjectedTableRecords(file *os.File) [][]types.UserFieldValue {
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

func AlterRelationRecord() {

}

func DeleteTableRecord(file *os.File, pkId int32, pkValue interface{}) error {
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
	var i1 int32
	for i, record := range allRecords {
		for _, fieldValue := range record {
			if fieldValue.ID == pkId && fieldValue.Value == pkValue {
				i1 = int32(i)
			}
		}
	}
	allRecords = append(allRecords[:i1], allRecords[i1+1:]...)

	file.Seek(0, 0) // Rewind to the beginning of the file
	err = json.NewEncoder(file).Encode(allRecords)
	if err != nil {
		return err
	}
	return nil
}
