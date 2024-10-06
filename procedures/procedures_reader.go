package procedures

import (
	"encoding/binary"
	"io"
	"myDb/types"
	"os"
)

func readFixedSizeString(file *os.File, stringLen int32) (string, error) {
	stringBytes := make([]byte, stringLen)
	if _, err := io.ReadFull(file, stringBytes); err != nil {
		return "", err
	}

	return string(stringBytes), nil
}

func readInt32(file *os.File) (int32, error) {
	var integer int32
	if err := binary.Read(file, binary.LittleEndian, &integer); err != nil {
		return -1, err
	}

	return integer, nil
}

func readFloat(file *os.File) (float64, error) {
	var float float64
	if err := binary.Read(file, binary.LittleEndian, &float); err != nil {
		return -1, err
	}

	return float, nil
}

func LoadTables(filename string) []types.Table {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	result := make([]types.Table, 0)
	var next int32
	for next != -1 {
		var table types.Table

		next, err = readInt32(file)
		if err == io.EOF {
			return make([]types.Table, 0)
		}
		table.Id, _ = readInt32(file)
		table.Size, _ = readInt32(file)
		nameLength, _ := readInt32(file)
		table.Name, _ = readFixedSizeString(file, nameLength)
		filenameLen, _ := readInt32(file)
		table.DataFileName, _ = readFixedSizeString(file, filenameLen)

		fieldsCount, _ := readInt32(file)
		table.Fields = make([]*types.Field, fieldsCount)
		for i := range table.Fields {
			table.Fields[i] = new(types.Field)
			table.Fields[i].FieldId, _ = readInt32(file)
			t, _ := readInt32(file)
			table.Fields[i].Type = types.DbType(t)
			table.Fields[i].Size, _ = readInt32(file)
			nameLen, _ := readInt32(file)
			table.Fields[i].Name, _ = readFixedSizeString(file, nameLen)
			table.Fields[i].Key, _ = readInt32(file)
		}
		result = append(result, table)
	}
	return result
}
