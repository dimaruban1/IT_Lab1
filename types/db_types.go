package types

// mysql types
type DbType rune

type FieldValue struct {
	ID        int32       `json:"id"`
	ValueType DbType      `json:"type"`
	Value     interface{} `json:"value"`
}

type ColorInvl struct {
	Color1          string
	Color2          string
	IntervalSeconds float64
}

type UserFieldValue struct {
	ID    int32       `json:"id"`
	Value interface{} `json:"value"`
}

type UserField struct {
	FieldId int32  `json:"id"`
	Type    string `json:"type"`
	Size    int32  `json:"size"`
	Name    string `json:"name"`
	Key     string `json:"key"`
}

type UserTable struct {
	Name   string      `json:"name"`
	Fields []UserField `json:"fields"`
}

const (
	// String types
	Char_t DbType = iota + 1000
	String_t

	// Numeric types
	Int_t
	Real_t

	// Date and time types
	Color_t
	ColorInvl_t

	notype_t
)

var DbTypes = [...]string{
	// String types
	"char",
	"string",

	// Numeric types
	"int",
	"real",

	// Date and time types
	"color",
	"colorinvl",
}

var SizeSpecifiedTypes = [...]string{
	"string",
}

var DbTypeDefaultSize = map[DbType]int32{
	Char_t:   1,
	String_t: 255,

	Int_t:  4,
	Real_t: 8,

	Color_t:     6,
	ColorInvl_t: 0,
}

var DbTypeMap = map[string]DbType{
	// String types
	"char":   Char_t,
	"string": String_t,

	// Numeric types
	"int":  Int_t,
	"real": Real_t,

	// Date and time types
	"color":     Color_t,
	"colorInvl": ColorInvl_t,
}

func ArrayContains(array []string, t string) bool {
	for _, a := range array {
		if a == t {
			return true
		}
	}
	return false
}

func GetDbTypeString(t DbType) string {
	for dbTypeString, dbTypeInt := range DbTypeMap {
		if dbTypeInt == t {
			return dbTypeString
		}
	}
	return ""
}
