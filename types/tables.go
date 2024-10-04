package types

import "fmt"

type KeyType = rune

const (
	PrimaryKey KeyType = 'P'
	ForeignKey KeyType = 'F'
	Nothing    KeyType = 'N'
)

type Field struct {
	FieldId int32   `json:"id"`
	Type    DbType  `json:"type"`
	Size    int32   `json:"size"`
	Name    string  `json:"name"`
	Key     KeyType `json:"key"`
}

type Table struct {
	Id           int32    `json:"id"`
	Size         int32    `json:"size"`
	Name         string   `json:"name"`
	Fields       []*Field `json:"fields"`
	DataFileName string   `json:"data_file_name"`
}

func NewField() *Field {
	field := new(Field)
	field.FieldId = 0
	field.Key = 'N'
	field.Name = ""
	field.Size = 0
	field.Type = notype_t
	return field
}

func (t *Table) ToString() string {
	result := ""

	result += fmt.Sprintf("\tid: %d\n", t.Id)
	result += fmt.Sprintf("\tsize: %d\n", t.Size)
	result += fmt.Sprintf("\tname: %s\n", t.Name)
	result += fmt.Sprintf("\tdata file name: %s\n", t.DataFileName)

	result += "\tfields:\n"
	for _, field := range t.Fields {
		result += field.ToString()
		result += "\n"
	}
	return result
}

func (field *Field) ToString() string {
	result := ""
	result += fmt.Sprintf("\t\tid: %d\n", field.FieldId)
	result += fmt.Sprintf("\t\ttype: %d\n", field.Type)
	result += fmt.Sprintf("\t\tsize: %d\n", field.Size)
	result += fmt.Sprintf("\t\tname: %s\n", field.Name)
	result += fmt.Sprintf("\t\tkey: %c\n", field.Key)

	return result
}

func (t *Table) GetFieldByName(fieldName string) *Field {
	for _, field := range t.Fields {
		if field.Name == fieldName {
			return field
		}
	}
	return nil
}

func (t *Table) GetFieldById(id int32) *Field {
	for _, field := range t.Fields {
		if field.FieldId == id {
			return field
		}
	}
	return nil
}
