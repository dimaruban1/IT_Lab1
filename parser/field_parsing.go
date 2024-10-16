package parser

import (
	"errors"
	"fmt"
	"myDb/types"
	"regexp"
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

func parseColorInvl(value string) *types.ColorInvl {

	return nil
}
