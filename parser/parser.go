package parser

import (
	"fmt"
	"myDb/types"
	"myDb/utility"
	"regexp"
	"strconv"
	"strings"
)

type Query struct {
	Text string
	Type QueryType
}

type QueryType int32

const (
	CreateQuery_t QueryType = iota
	InsertRecordQuery_t
	InsertDatasetQuery_t
	Misc_t //Error type
)

const (
	createQueryRegex = iota
	insertRecordQueryRegex
	insertDatasetRegex
	identifierRegex
	bracketRegex
	tokenRegex
	greedyTokenRegex
	sizeOfTypeRegex
	fieldValueRegex
)

var actions = `(CREATE|ALTER|DELETE)`
var objects = `(RELATION|DATASET)`

var regexMap = map[int]string{
	createQueryRegex:       `(?im)` + actions + `(\s+)` + objects + `(\s+)[a-zA-Z]\w*(\s+)\((?s).*\)`,
	insertRecordQueryRegex: `INSERT\s+INTO\s+\w+\s*\(.+?\)\s+VALUES\s*\(.+?\)`,
	insertDatasetRegex:     `INSERT\s+INTO\s+\w+\s+OWNER\s*\(\w+\)\s+MEMBER\s*\(\w+(,\s*\w+)*\)`,
	identifierRegex:        `[a-zA-Z]\w*`,
	bracketRegex:           `\((?s).*\)`,
	tokenRegex:             `[a-zA-Z]\w*|\((?s).*?\)`,
	greedyTokenRegex:       `[a-zA-Z]\w*|\((?s).*\)`,
	sizeOfTypeRegex:        `\([1-9]\d*\)`,
	fieldValueRegex:        `\".*?\"`,
}

func isRegexCorrect(stringLiteral, regex string) bool {
	strings, err := regexThatString(stringLiteral, regex)
	if err != nil {
		return false
	}
	return len(strings) == 1
}

func isQueryCorrect(query Query) bool {
	strings, err := getStringsOfRegex(query.Text, int(query.Type))
	if err != nil {
		return false
	}
	return len(strings) == 1
}

func regexThatString(stringLiteral, regex string) ([]string, error) {
	myRegex := regexp.MustCompile(regex)
	matchedResults := myRegex.FindAllStringSubmatch(stringLiteral, -1)

	var strings []string
	for _, match := range matchedResults {
		strings = append(strings, match[0])
	}

	return strings, nil
}

func getStringsOfRegex(stringLiteral string, regexType int) ([]string, error) {
	regex, ok := regexMap[regexType]
	if !ok {
		return nil, fmt.Errorf("regex type %d doesn`t exist", regexType)
	}
	myRegex := regexp.MustCompile(regex)
	matchedResults := myRegex.FindAllStringSubmatch(stringLiteral, -1)

	var strings []string
	for _, match := range matchedResults {
		strings = append(strings, match[0])
	}

	return strings, nil
}

func ParseFieldValue(fieldValue *types.FieldValue, value string) error {
	switch fieldValue.ValueType {
	case types.Int_t:
		parsedValue, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		fieldValue.Value = parsedValue
	case types.Real_t:
		parsedValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		fieldValue.Value = parsedValue
	case types.Char_t, types.String_t:
		fieldValue.Value = value
	case types.Color_t:
		parsedValue, err := ParseColor(value)
		if err != nil {
			return err
		}
		fieldValue.Value = parsedValue
	case types.ColorInvl_t:
		ParseColorInvl(value)
	default:
		return fmt.Errorf("unsupported ValueType: %d", fieldValue.ValueType)
	}
	return nil
}

func ParseColorInvl(value string) {

}

func ParseInsertRecordQuery(insertQuery string) (string, []map[string]string, error) {
	if !isQueryCorrect(Query{Text: insertQuery, Type: InsertRecordQuery_t}) {
		return "", nil, fmt.Errorf("query '%s' is incorrect", insertQuery)
	}
	tokens, err := getStringsOfRegex(insertQuery, tokenRegex)
	if err != nil || len(tokens) < 6 {
		return "", nil, fmt.Errorf("error when parsing query '%s'", insertQuery)
	}

	// INSERT INTO <tablename> => <tablename> index 2
	tableName := tokens[2]

	fieldsBracket := tokens[3]
	removeBrackets(&fieldsBracket)

	fieldNames := strings.Split(utility.RemoveWhitespaces(fieldsBracket), ",")

	var tuples []map[string]string
	for i := 5; i < len(tokens); i++ {
		tuple, err := parseInsertTupleBrackets(fieldNames, tokens[i])
		if err != nil {
			return "", nil, err
		}
		tuples = append(tuples, tuple)
	}
	return tableName, tuples, nil
}

func ParseCreateTableQuery(createQuery string) (*types.Table, error) {
	if !isQueryCorrect(Query{Text: createQuery, Type: CreateQuery_t}) {
		return nil, fmt.Errorf("query '%s' is incorrect", createQuery)
	}
	table := new(types.Table)
	table.Fields = make([]*types.Field, 0)
	tokens, err := getStringsOfRegex(createQuery, greedyTokenRegex)
	if err != nil || len(tokens) < 4 {
		return nil, fmt.Errorf("error when parsing query '%s'", createQuery)
	}

	table.Name = tokens[2]
	fieldTokens, err := getFieldsFromQuery(createQuery)
	if err != nil {
		return nil, err
	}

	var size int32 = 0
	for i, fieldToken := range fieldTokens {
		field, err := parseFieldToken(fieldToken, int32(i))
		if err != nil {
			return nil, err
		}
		table.Fields = append(table.Fields, field)
		size += field.Size
	}
	table.RecordsCount = 0
	table.DataFileName = table.Name + "_table.bin"
	table.Size = size

	return table, nil
}

func parseInsertTupleBrackets(fieldNames []string, tupleBracketString string) (map[string]string, error) {
	// tupleBracketString = utility.RemoveWhitespaces(tupleBracketString)
	values, err := getStringsOfRegex(tupleBracketString, fieldValueRegex)
	if err != nil {
		return nil, err
	}

	if len(values) != len(fieldNames) {
		return nil,
			fmt.Errorf("expected number of fields %d, got number of fields %d in string (%s)",
				len(values), len(fieldNames), tupleBracketString)
	}

	var res map[string]string = make(map[string]string)
	for i, v := range values {
		removeBrackets(&v)
		res[fieldNames[i]] = v
	}
	return res, nil
}

func ParseAlterTableQuery() {

}
