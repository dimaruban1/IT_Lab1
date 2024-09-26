package utility

import (
	"strings"
	"testing"
)

func TestFlattenWhitespace(t *testing.T) {
	str := `CREATE RELATION Persons (
		PersonID int,
		LastName varchar(255),
		FirstName varchar(255),
		Address varchar(255),
		City varchar(255)
	);`
	res := FlattenWhitespaces(str)

	if strings.Contains(res, "\n") || strings.Contains(res, "\t") {
		t.Error(res)
	}
}