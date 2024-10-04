package recording

import (
	"io"
	"myDb/params"
	"myDb/procedures"
	SysCatalog "myDb/system_catalog"
	"myDb/types"
	"os"
)

func GetRecords(tableName string) []map[types.Field]*types.FieldValue {
	table := SysCatalog.GetTableByName(tableName)
	if table == nil {
		return nil
	}
	filename := table.DataFileName
	records := make([]map[types.Field]*types.FieldValue, 0)

	file, err := os.Open(params.SaveDir + "\\" + filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for {
		record, err := procedures.ReadRecord(file, table)

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		records = append(records, record)
	}
	return records
}

func ProjectRecords(tableName string, fieldNames []string) [][]*types.FieldValue {
	table := SysCatalog.GetTableByName(tableName)
	filename := table.DataFileName
	projectedRecords := make([][]*types.FieldValue, 0)

	fieldIds := make([]int32, len(fieldNames))
	for _, fieldName := range fieldNames {
		field := table.GetFieldByName(fieldName)
		if field == nil {
			// error
		}
		fieldIds = append(fieldIds, field.FieldId)
	}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for {
		record, err := procedures.ReadRecord(file, table)
		projectedRecord := make([]*types.FieldValue, len(fieldNames))

		for _, fieldValue := range record {
			for _, fid := range fieldIds {
				if fieldValue.ID == fid {
					projectedRecord = append(projectedRecord, fieldValue)
				}

			}
		}

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		projectedRecords = append(projectedRecords, projectedRecord)
	}
	return projectedRecords
}
