package types

// mysql types
type DbType rune

type ColorInvl_struct struct {
	color1           FieldValue
	color2           FieldValue
	interval_seconds int32
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

type FieldValue struct {
	ID        int32
	ValueType DbType
	Value     interface{}
}

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

	Color_t:     0,
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
