package parser

import (
	"fmt"
	"myDb/types"
	"regexp"
	"strconv"
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

func isRegexCorrect(stringLiteral, regex string) bool {
	strings, err := regexThatString(stringLiteral, regex)
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

func ParseFieldValue(field *types.Field, value string) (interface{}, error) {
	switch field.Type {
	case types.Int_t:
		parsedValue, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		return parsedValue, nil
	case types.Real_t:
		parsedValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return parsedValue, nil
	case types.Char_t, types.String_t:
		return value, nil
	case types.Color_t:
		parsedValue, err := ParseColor(value)
		if err != nil {
			return nil, err
		}
		return parsedValue, nil
	case types.ColorInvl_t:
		ParseColorInvl(value)
	default:
		return nil, fmt.Errorf("unsupported ValueType: %d", field.Type)
	}
	return nil, nil
}

// colorInvl string:
func ParseColorInvl(value string) {

}
