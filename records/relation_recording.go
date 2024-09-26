package recording

import (
	"io"
	"myDb/procedures"
	"myDb/types"
	"os"
)

func WriteTableRecord(file *os.File, objectFieldValues []types.FieldValue, offset64 int64) {
	if offset64 >= 0 {
		file.Seek(offset64, 0)
	} else {
		file.Seek(0, io.SeekEnd)
	}
	for _, field := range objectFieldValues {
		procedures.WriteField(field, file)
	}
}

func InsertRelationRecord() {

}

func AlterRelationRecord() {

}

func DeleteRelationRecord() {

}
