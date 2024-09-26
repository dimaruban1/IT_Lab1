package recording

import (
	"io"
	"myDb/procedures"
	SysCatalog "myDb/system_catalog"
	"myDb/types"
	"os"
)

func GetRecords(tableName string) [][]*types.FieldValue {
	table := SysCatalog.GetTableByName(tableName)
	filename := table.DataFileName
	records := make([][]*types.FieldValue, 0)

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for {
		record, err := readRecord(file, table)

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
		field := table.GetField(fieldName)
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
		record, err := readRecord(file, table)
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

func readRecord(file *os.File, table *types.Table) ([]*types.FieldValue, error) {
	record := make([]*types.FieldValue, len(table.Fields))
	for _, field := range table.Fields {
		fieldValue, err := procedures.ReadField(*field, file)
		if err != nil {
			return nil, err
		}
		record = append(record, fieldValue)
	}
	return record, nil
}
