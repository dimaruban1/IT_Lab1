package SysCatalog

import (
	"fmt"
	"myDb/types"
)

// add editing
// add deleting
// add colorInvl
// test color
// finish GREAT LOOP
var Tables []types.Table
var CurrentRecords []map[types.Field]*types.FieldValue

func NewDB() {
	Tables = make([]types.Table, 0)
	CurrentRecords = make([]map[types.Field]*types.FieldValue, 0)
}

func GetTableByName(name string) *types.Table {
	for _, table := range Tables {
		if table.Name == name {
			return &table
		}
	}
	return nil
}

func GetTableById(id int32) *types.Table {
	for _, table := range Tables {
		if table.Id == id {
			return &table
		}
	}
	return nil
}

func DeleteTableByName(name string) error {
	for j, relation := range Tables {
		if relation.Name == name {
			Tables = append(Tables[:j], Tables[j+1:]...)
			return nil
		}
	}
	return fmt.Errorf("таблицю %s не видалено, помилка: таблицю не знайдено", name)
}
