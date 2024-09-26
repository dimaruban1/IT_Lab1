package procedures

import (
	"encoding/binary"
	"fmt"
	"io"
	"myDb/types"
	"myDb/utility"
	"os"
)

func SaveAllTablesBin(tables []types.Table, filename string) {
	utility.CreateFileIfNotExists(filename)
	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// store with smallest first
	// sort.Sort(SysCatalog.RelationListSort(tables))

	var offset int32 = 0
	for i, table := range tables {
		isLast := i == len(tables)-1
		offset = WriteTableToFile(file, offset, table, isLast)
	}
}

func InsertTable(tuple map[int]string, table *types.Table) {
	filename := table.DataFileName
	utility.CreateFileIfNotExists(filename)

	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	file.Seek(0, io.SeekEnd)

	for fieldId, fieldValue := range tuple {
		binary.Write(file, binary.LittleEndian, int32(fieldId))
		binary.Write(file, binary.LittleEndian, int32(len(fieldValue)))
		binary.Write(file, binary.LittleEndian, []byte(fieldValue))
	}

}

func WriteTableToFile(file *os.File, offset int32, table types.Table, isLast bool) int32 {
	offset64 := int64(offset)
	_, err := file.Seek(offset64, 0)
	if err != nil {
		panic(err)
	}

	binary.Write(file, binary.LittleEndian, offset)
	binary.Write(file, binary.LittleEndian, table.Id)
	binary.Write(file, binary.LittleEndian, int32(table.Size))
	binary.Write(file, binary.LittleEndian, int32(len(table.Name)))
	binary.Write(file, binary.LittleEndian, []byte(table.Name))
	binary.Write(file, binary.LittleEndian, int32(len(table.DataFileName)))
	binary.Write(file, binary.LittleEndian, []byte(table.DataFileName))
	binary.Write(file, binary.LittleEndian, int32(table.RecordsCount))

	binary.Write(file, binary.LittleEndian, int32(len(table.Fields)))
	for _, field := range table.Fields {
		binary.Write(file, binary.LittleEndian, int32(field.FieldId))
		binary.Write(file, binary.LittleEndian, field.Type)
		binary.Write(file, binary.LittleEndian, int32(field.Size))
		binary.Write(file, binary.LittleEndian, int32(len(field.Name)))
		binary.Write(file, binary.LittleEndian, []byte(field.Name))
		binary.Write(file, binary.LittleEndian, field.Key)
	}

	length, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		panic(err)
	}

	newOffset := int32(length)
	if isLast {
		newOffset = -1
	}
	file.Seek(offset64, 0)
	binary.Write(file, binary.LittleEndian, newOffset)
	return newOffset
}

func WriteField(field types.FieldValue, file *os.File) {
	switch field.ValueType {
	case types.Char_t:
	case types.String_t:
	case types.Color_t:
		if v, ok := field.Value.(string); ok {
			binary.Write(file, binary.LittleEndian, []byte(v))
		} else {
			fmt.Printf("Wrong assumed format error")
			break
		}

	case types.Int_t:
		if v, ok := field.Value.(int); ok {
			err := binary.Write(file, binary.LittleEndian, int32(v))
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("Wrong assumed format error for floating-point types")
			break
		}

	case types.Real_t:
		if v, ok := field.Value.(float64); ok {
			if err := binary.Write(file, binary.LittleEndian, v); err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("Wrong assumed format error for floating-point types")
			break
		}

	// TODO: implement
	case types.ColorInvl_t:

	}
}
