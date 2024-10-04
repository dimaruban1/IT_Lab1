package types

func CastToUserTable(table Table) *UserTable {
	var userTable = new(UserTable)
	userTable.Name = table.Name
	userTable.Fields = make([]UserField, 0)

	for _, field := range table.Fields {
		f := CastToUserField(*field)
		userTable.Fields = append(userTable.Fields, *f)
	}

	return userTable
}

func CastFromUserTable(id int, userTable UserTable) *Table {
	var table = new(Table)
	table.Size = 0
	table.Fields = make([]*Field, 0)
	table.Id = int32(id)
	table.Name = userTable.Name

	for _, field := range userTable.Fields {
		f := CastFromUserField(field)
		table.Fields = append(table.Fields, f)
		table.Size += field.Size
	}

	table.DataFileName = table.Name + "_table.json"

	return table
}

func CastToUserField(field Field) *UserField {
	userField := new(UserField)
	userField.Name = field.Name
	userField.FieldId = field.FieldId
	userField.Key = string(field.Key)
	userField.Size = field.Size
	userField.Type = GetDbTypeString(field.Type)

	return userField
}

func CastFromUserField(field UserField) *Field {
	f := new(Field)
	f.FieldId = field.FieldId
	f.Key = rune(field.Key[0])
	f.Type = DbTypeMap[field.Type]
	f.Name = field.Name
	f.Size = field.Size

	return f
}

func CastToUserFieldValue(fieldValue FieldValue) *UserFieldValue {
	ufv := new(UserFieldValue)
	ufv.ID = fieldValue.ID
	ufv.Value = fieldValue.Value

	return ufv
}

func CastFromUserFieldValue(field Field, userFieldValue UserFieldValue) *FieldValue {
	fv := new(FieldValue)
	fv.ID = userFieldValue.ID
	fv.Value = userFieldValue.Value
	fv.ValueType = field.Type

	return fv
}
