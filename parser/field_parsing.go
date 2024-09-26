package parser

import (
	"errors"
	"fmt"
	"myDb/types"
	"myDb/utility"
	"regexp"
	"strconv"
	"strings"
)

// default one is 'AABBCC'
var colorRegexes = [...]string{
	`#([a-fA-F0-9]{6}|[a-fA-F0-9]{3})`,
	`([a-fA-F0-9]{6}|[a-fA-F0-9]{3})`,
	`rgb\(\s*\d{1,3}\s*,\s*\d{1,3}\s*,\s*\d{1,3}\s*\)/gm`,

	// implement later
	`\b(red|blue|green|yellow|black|white|gray)\b`,
}

func expandShortHex(shortHex string) string {
	return fmt.Sprintf("#%c%c%c%c%c%c", shortHex[0], shortHex[0], shortHex[1], shortHex[1], shortHex[2], shortHex[2])
}

// Convert RGB format to hex
func rgbToHex(rgb string) (string, error) {
	// Find all integers in the rgb string
	re := regexp.MustCompile(`\d{1,3}`)
	values := re.FindAllString(rgb, -1)
	if len(values) != 3 {
		return "", errors.New("invalid rgb format")
	}

	var r, g, b int
	fmt.Sscanf(values[0], "%d", &r)
	fmt.Sscanf(values[1], "%d", &g)
	fmt.Sscanf(values[2], "%d", &b)

	return fmt.Sprintf("#%02x%02x%02x", r, g, b), nil
}

func ParseColor(value string) (string, error) {
	for i, colorRegex := range colorRegexes {
		if isRegexCorrect(value, colorRegex) {
			switch i {
			// Case 0: #AABBCC or #ABC hex format
			case 0:
				if len(value) == 4 {
					// Short hex, expand it
					return expandShortHex(value[1:]), nil
				}
				return strings.ToUpper(value)[1:], nil

			// Case 1: AABBCC or ABC without #
			case 1:
				if len(value) == 3 {
					return expandShortHex(value), nil
				}
				return strings.ToUpper(value), nil

			// Case 2: rgb(r, g, b) format
			case 2:
				hex, err := rgbToHex(value)
				if err != nil {
					return "", err
				}
				return strings.ToUpper(hex), nil

			// Case 3: Color names (red, blue, etc.) - can be handled later
			case 3:
				// Implement named colors mapping later
				return "", fmt.Errorf("color names are not implemented yet")
			}
		}
	}
	return "", fmt.Errorf("color string '%s' incorrect", value)
}

func parseColorInvl(value string) (string, error) {
	return "", fmt.Errorf("no implement")
}

func parseFieldToken(fieldToken string, fieldId int32) (*types.Field, error) {
	field := new(types.Field)
	field.FieldId = fieldId

	tokens := strings.Split(fieldToken, " ")
	if tokens[0] == "" {
		tokens[0] = tokens[1]
		tokens[1] = tokens[2]
	}
	field.Name = tokens[0]
	typeString := tokens[1]
	typeName, err := getStringsOfRegex(typeString, identifierRegex)
	typeName0 := typeName[0]
	if err != nil {
		return nil, err
	}
	if field.Name[len(field.Name)-1] == '#' {
		field.Key = 'P'
	} else {
		field.Key = 'N'
	}

	if !strings.Contains(fieldToken, "(") && types.ArrayContains(types.DbTypes[:], typeName0) {
		field.Type = types.DbTypeMap[typeName0]
		field.Size = types.DbTypeDefaultSize[field.Type]

		return field, nil
	}

	sizeBraket, _ := getStringsOfRegex(typeString, bracketRegex)

	if strings.Contains(typeName0, "(") && !types.ArrayContains(types.SizeSpecifiedTypes[:], typeName0) {
		return nil, fmt.Errorf("field %s is not size specified", typeName0)
	}
	if !types.ArrayContains(types.DbTypes[:], typeName0) {
		return nil, fmt.Errorf("failed to parse fields, unknown type: %s", typeName0)
	}

	removeBrackets(&sizeBraket[0])
	sizeString := sizeBraket[0]
	size, err := strconv.Atoi(sizeString)
	if err != nil {
		return nil, err
	}

	field.Type = types.DbTypeMap[typeName0]
	field.Size = int32(size)

	return field, nil
}

func removeBrackets(token *string) {
	if len(*token) >= 2 &&
		((*token)[0] == '(' && (*token)[len(*token)-1] == ')' ||
			(*token)[0] == '"' && (*token)[len(*token)-1] == '"') {
		*token = (*token)[1 : len(*token)-1]
	}
}

func getFieldsFromQuery(fullQuery string) ([]string, error) {
	bracketTokens, err := getStringsOfRegex(fullQuery, bracketRegex)
	if err != nil {
		return nil, err
	}
	if len(bracketTokens) == 0 {
		return nil, errors.New("no brackets found")
	}

	bracketToken := bracketTokens[0]
	removeBrackets(&bracketToken)

	bracketToken = utility.FlattenWhitespaces(bracketToken)
	fieldTokens := strings.Split(bracketToken, ",")

	return fieldTokens, nil
}

func ProcessInsertion(values []map[string]string, table *types.Table) ([][]types.FieldValue, error) {
	result := make([][]types.FieldValue, len(values))
	for i := range result {
		result[i] = make([]types.FieldValue, len(table.Fields))
	}

	fieldIndexMap := make(map[string]int)
	for i, field := range table.Fields {
		fieldIndexMap[field.Name] = i
	}

	for rowIndex, valueEntry := range values {
		for nameKey, value := range valueEntry {
			fieldIndex, ok := fieldIndexMap[nameKey]
			if !ok {
				continue
			}

			field := table.Fields[fieldIndex]

			result[rowIndex][fieldIndex].ID = field.FieldId
			result[rowIndex][fieldIndex].ValueType = types.DbType(field.Type)

			err := ParseFieldValue(&result[rowIndex][fieldIndex], value)
			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}
