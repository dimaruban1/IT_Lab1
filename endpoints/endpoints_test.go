package endpoints

import (
	"myDb/procedures"
	SysCatalog "myDb/system_catalog"
	"myDb/types"
	"os"
	"testing"
)

func getBooksTable() *types.Table {
	fields := []*types.Field{
		{FieldId: 1, Type: 1, Size: 255, Name: "title", Key: 'P'},
		{FieldId: 2, Type: 1, Size: 255, Name: "author", Key: 'N'},
		{FieldId: 3, Type: 2, Size: 4, Name: "year_published", Key: 'N'},
		{FieldId: 4, Type: 1, Size: 255, Name: "isbn", Key: 'N'},
	}

	// Create the books table
	booksTable := &types.Table{
		Id:           1,
		Size:         int32(len(fields)),
		Name:         "books",
		Fields:       fields,
		DataFileName: "books_data.json",
	}

	return booksTable
}

var testFilename = "test_tables.bin"
var booksTableTestFileName = "test_books_data.json"

func cleanUp() {
	os.Remove(testFilename)
	os.Remove(booksTableTestFileName)
}

// func TestCreateRecord(t *testing.T) {
// 	table := getBooksTable()

// 	utility.CreateFileIfNotExists(booksTableTestFileName)
// 	file, err := os.OpenFile(booksTableTestFileName, os.O_RDWR, 0666)
// 	if err != nil {
// 		t.Fail()
// 		cleanUp()
// 	}

// 	// Test valid inputs
// 	// fieldValues := map[int32]interface{}{
// 	// 	1: 12,
// 	// 	2: "greybeards",
// 	// 	3: "FFFFFF",
// 	// 	4: 12.56,
// 	// }

// 	// r := createRecord(table, fieldValues, file)
// 	if len(r) != 4 {
// 		t.Fail()
// 		cleanUp()
// 	}
// }

func TestGetRecords(t *testing.T) {
	table := getBooksTable()

	fieldValues := getRecords(table, booksTableTestFileName)
	if len(fieldValues) == 0 {
		t.Fail()
		cleanUp()
	}
	if len(fieldValues) != 1 || len(fieldValues[0]) != 4 {
		t.Fail()
		cleanUp()
	}
}

func TestAlterRecord(t *testing.T) {
	table := getBooksTable()

	// Test valid inputs
	recordIdToDelete := 12

	newFieldValues := map[int32]interface{}{
		1: 12,
		2: "whitebeards",
		3: "SSSSSS",
		4: 12.56,
	}
	_ = alterRecord(table, booksTableTestFileName, recordIdToDelete, newFieldValues)

	tables := getRecords(table, booksTableTestFileName)

	if len(tables) != 1 {
		t.Fail()
		cleanUp()
	}
	if tables[0][2] != "whitebeards" {
		t.Fail()
		cleanUp()
	}
}

func TestDeleteRecord(t *testing.T) {
	table := getBooksTable()

	// Test valid inputs
	recordIdToDelete1 := 12
	recordIdToDelete2 := 1

	_ = deleteRecord(table, booksTableTestFileName, recordIdToDelete2)
	tables := getRecords(table, booksTableTestFileName)

	// if getRecord(table, booksTableTestFileName, 12) == nil {
	// 	t.Fail()
	// }

	if len(tables) != 1 {
		t.Fail()
	}

	_ = deleteRecord(table, booksTableTestFileName, recordIdToDelete1)

	tables = getRecords(table, booksTableTestFileName)

	if len(tables) != 0 {
		t.Fail()
	}

	cleanUp()
}

func TestDeleteRecord2(t *testing.T) {
	SysCatalog.Tables = procedures.LoadTables("C:\\Users\\Dima\\go\\inf_technologies_lab_1\\data\\tables.bin")
	table := SysCatalog.GetTableByName("books")
	file, _ := os.Open("C:\\Users\\Dima\\go\\inf_technologies_lab_1\\data\\books_table.json")
	t1 := getRecord(table, file, "1")
	t2 := getRecord(table, file, 1)
	if len(t1) != len(t2) {
		t.Fail()
	}
}

func TestProjectRecords(t *testing.T) {
	SysCatalog.Tables = procedures.LoadTables("C:\\Users\\Dima\\go\\inf_technologies_lab_1\\data\\tables.bin")
	table := SysCatalog.GetTableByName("books")
	file, _ := os.Open("C:\\Users\\Dima\\go\\inf_technologies_lab_1\\data\\books_table.json")

	fields := make([]types.Field, 0)
	fields = append(fields, *table.Fields[0])
	fields = append(fields, *table.Fields[2])

	projected_records := getProjectedTableRecords(file, fields)
	if len(projected_records[0]) != 2 {
		t.Fail()
	}
}
