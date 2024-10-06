package endpoints

import (
	"myDb/procedures"
	SysCatalog "myDb/system_catalog"
	"myDb/types"
	"myDb/utility"
	"os"
	"testing"
)

func cleanUp() {

}

func getBooksTable() *types.Table {
	fields := []*types.Field{
		{FieldId: 1, Type: 1, Size: 255, Name: "title", Key: 0},
		{FieldId: 2, Type: 1, Size: 255, Name: "author", Key: 0},
		{FieldId: 3, Type: 2, Size: 4, Name: "year_published", Key: 0},
		{FieldId: 4, Type: 1, Size: 255, Name: "isbn", Key: 1},
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

func TestCreateRecord(t *testing.T) {
	table := getBooksTable()

	utility.CreateFileIfNotExists(booksTableTestFileName)
	file, err := os.Open(booksTableTestFileName)
	if err != nil {
		t.Fail()
		cleanUp()
	}

	// Test valid inputs
	fieldValues := map[int32]interface{}{
		1: 12,
		2: "greybeards",
		3: "FFFFFF",
		4: 12.56,
	}

	r := createRecord(table, fieldValues, file)
	if len(r) != 4 {
		t.Fail()
	}
}

func TestGetRecords(t *testing.T) {
	SysCatalog.Tables = procedures.LoadTables("C:\\Users\\Dima\\go\\inf_technologies_lab_1\\data\\tables.bin")
	// Test valid inputs
	table := SysCatalog.GetTableByName("books")
	fieldValues := getRecords(table, "C:\\Users\\Dima\\go\\inf_technologies_lab_1\\data\\books_table.json")
	if len(fieldValues) == 0 {
		t.Fail()
	}
	if len(fieldValues) != 3 {
		t.Fail()
	}
}

func TestDeleteRecord(t *testing.T) {
	SysCatalog.Tables = procedures.LoadTables("C:\\Users\\Dima\\go\\inf_technologies_lab_1\\data\\tables.bin")
	table := SysCatalog.GetTableByName("books")
	utility.CreateFileIfNotExists("C:\\Users\\Dima\\go\\inf_technologies_lab_1\\data\\books_table.json")
	file, err := os.OpenFile("C:\\Users\\Dima\\go\\inf_technologies_lab_1\\data\\books_table.json", os.O_RDWR, 0666)
	if err != nil {
		t.Fail()
	}

	// Test valid inputs
	fieldValues := map[int32]interface{}{
		1: 12,
		2: "greybeards",
		3: "FFFFFF",
		4: 12.56,
	}

	_ = createRecord(table, fieldValues, file)
}
