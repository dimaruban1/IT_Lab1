package CLI

import "strings"

type Command struct {
	Usage       string
	minArgCount int
}

const DefaultDatasetFilename = "system_catalog\\datasets.bin"
const DefaultRelationFilename = "system_catalog\\relations.bin"

var Commands = map[string]Command{
	"save": {
		Usage:       "SAVE TABLES [<FILENAME>]",
		minArgCount: 2,
	},
	"print": {
		Usage:       "PRINT TABLES",
		minArgCount: 2,
	},
	"load": {
		Usage:       "LOAD TABLES [<FILENAME>]",
		minArgCount: 2,
	},
	"create": {
		Usage:       "CREATE TABLE <FILENAME>",
		minArgCount: 3,
	},
	"delete": {
		Usage:       "DELETE TABLE <OBJECT NAME>",
		minArgCount: 3,
	},
	"list": {
		Usage:       "LIST",
		minArgCount: 1,
	},
	"set": {
		Usage:       "SET WORKDIR|SAVEDIR <PATH>",
		minArgCount: 3,
	},
	"savedir": {
		Usage:       "SAVEDIR",
		minArgCount: 1,
	},
	"workdir": {
		Usage:       "WORKDIR",
		minArgCount: 1,
	},
	"insert": {
		Usage:       "INSERT TABLE <FILENAME>",
		minArgCount: 3,
	},
}

func IsUsageCorrect(userInput string) bool {
	args := strings.Split(userInput, " ")
	command := strings.ToLower(args[0])
	return len(args) >= Commands[command].minArgCount
}

func GetArgumentsFromCommand(userInput string) (string, string, string) {
	tokens := strings.Split(userInput, " ")
	command := tokens[0]

	if len(tokens) < 2 {
		return command, "", ""
	}

	object := tokens[1]

	var filename string
	if object == "datasets" || object == "dataset" {
		filename = DefaultDatasetFilename
	} else if object == "relation" || object == "relations" {
		filename = DefaultRelationFilename
	} else {
		filename = ""
	}

	if len(tokens) > 2 {
		filename = tokens[2]
	}

	return command, object, filename
}

func CommandExists(command string) bool {
	for key := range Commands {
		if key == command {
			return true
		}
	}
	return false
}
