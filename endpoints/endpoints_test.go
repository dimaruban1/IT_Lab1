package endpoints

import (
	"myDb/procedures"
	SysCatalog "myDb/system_catalog"
	"myDb/utility"
	"os"
	"testing"
)

func TestCreateTable_Success(t *testing.T) {
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
