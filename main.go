package main

/*
CREATE table create_table\create_books.txt

*/

import (
	"bufio"
	"fmt"
	CLI "myDb/command_line_interface"
	"myDb/params"
	"myDb/parser"
	"myDb/procedures"
	recording "myDb/records"
	SysCatalog "myDb/system_catalog"
	"myDb/types"
	"myDb/utility"
	"os"
	"strings"
)

func main() {
	// process some commands
	// listCommands()
	// launchProgram()

	// cration of table and saving to file
	// createTable("example_queries\\create_table\\create_books.txt")
	// fmt.Printf("Кількість таблиць: %d\n", len(SysCatalog.Tables))
	// for _, table := range SysCatalog.Tables {
	// 	fmt.Print(table.ToString())
	// 	fmt.Print("\n\n\n")
	// }

	// filename := params.SaveDir + "\\" + "tables.bin"
	// procedures.SaveAllTablesBin(SysCatalog.Tables, filename)

	filename := params.SaveDir + "\\" + "tables.bin"
	SysCatalog.Tables = procedures.LoadTables(filename)
	filename = params.WorkDir + "\\insert_record\\insert_book.txt"
	//insertRecord(filename)

	printTables()
}

func acceptUserInput(message string) (string, error) {
	fmt.Print(message)
	var line string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line = scanner.Text()
	}

	return line, nil
}

func launchProgram() {
	for {
		input, err := acceptUserInput(">")
		if err != nil {
			fmt.Printf("Error encountered: %s", err.Error())
			continue
		}

		command, object, filename := CLI.GetArgumentsFromCommand(input)
		command = strings.ToLower(command)
		object = strings.ToLower(object)
		if !CLI.CommandExists(command) {
			fmt.Printf("Команди %s не існує, спробуйте одну з цих:\n", command)
			listCommands()
		}
		if !CLI.IsUsageCorrect(input) {
			fmt.Printf("Некоректно використано команду %s, правильно так:\n %s\n", command, CLI.Commands[command].Usage)
			continue
		}

		switch command {
		case "list":
			fmt.Println("Selected: List items")
			listCommands()

		case "create":
			filename = params.WorkDir + "\\" + filename
			switch object {
			case "table":
				createTable(filename)
			}
			fmt.Printf("\nВибрано: Виконати %s запит для %s файлу зі шляхом '%s'\n", command, object, filename)

		case "save":
			filename = params.SaveDir + "\\" + filename
			switch object {
			case "tables":
				procedures.SaveAllTablesBin(SysCatalog.Tables, filename)
			}
			fmt.Printf("\nВибрано: Виконати %s запит для %s файлу зі шляхом '%s'\n", command, object, filename)

		case "print":
			switch object {
			case "tables":
				if len(SysCatalog.Tables) == 0 {
					fmt.Println("Жодної таблиці ще не створено, неможливо презентувати")
					continue
				}
				printTables()
			}
			fmt.Printf("\nВибрано: Виконати %s запит для %s \n", command, object)
		case "load":
			filename = params.SaveDir + "\\" + filename
			switch object {
			case "tables":
				SysCatalog.Tables = procedures.LoadTables(filename)
			}
			fmt.Printf("\nВибрано: Виконати %s запит для %s файлу зі шляхом '%s'\n", command, object, filename)

		case "insert":
			filename = params.WorkDir + "\\" + filename
			switch object {
			case "table":
				err := insertRecord(filename)
				if err == nil {
					fmt.Print("Запис додано успішно")
				} else {
					fmt.Printf("Сталася помилка: %s", err.Error())
				}
			}

		case "delete":
			switch object {
			case "table":
				deleteTable(filename)
			}

		case "set":
			switch object {
			case "workdir":
				setWorkdir(filename)
			case "savedir":
				setSaveDir(filename)
			}
			fmt.Printf("\nВибрано: Виконати %s запит для %s файлу зі шляхом '%s'\n", command, object, filename)

		case "savedir":
			fmt.Printf("savedir: %s", params.SaveDir)

		case "workdir":
			fmt.Printf("workdir: %s", params.WorkDir)

		case "exit":
			fmt.Println("Selected: Exit the program")
			os.Exit(0)

		default:
			fmt.Println("Некоректна команда, спробуйте іншу")
		}
		fmt.Println()
	}
}

func listCommands() {
	fmt.Println("LIST. Список команд")
	fmt.Println("SAVE TABLES {FILENAME}. зберегти об'єкти у вказаний файл")
	fmt.Println("PRINT TABLES. надрукувати всі відношення")
	fmt.Println("LOAD TABLES {FILENAME}. завантажити всі набори даних з файлу")
	fmt.Println("CREATE TABLE {FILENAME}. створити набір даних або відношення")
	fmt.Println("SET SAVEDIR|WORKDIR {PATH}. встановити директорію для збереження або робочу директорію")
	fmt.Println("SAVEDIR|WORKDIR. показати відповідну встановлену директорію")
	fmt.Println("INSERT table {FILENAME}. вставити відношення з файлу")
	fmt.Println("INSERT DATASET {FILENAME}. вставити набір даних з файлу")
	fmt.Println("EXIT. Вихід з програми")
}

func createTable(filename string) {
	query, err := os.ReadFile(filename)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	elem, err := parser.ParseCreateTableQuery(string(query))
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	name := elem.Name
	table := SysCatalog.GetTableByName(name)
	if table != nil {
		fmt.Printf("Таблиця '%s' уже існує", name)
		return
	}
	SysCatalog.Tables = append(SysCatalog.Tables, *elem)
	fmt.Printf("таблицю %s успішно створено, нова кількість таблиць - %d\n", name, len(SysCatalog.Tables))
}

func insertRecord(filename string) error {
	query, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	tableName, fieldValues, err := parser.ParseInsertRecordQuery(string(query))
	if err != nil {
		return err
	}
	table := SysCatalog.GetTableByName(tableName)
	tuples, err := parser.ProcessInsertion(fieldValues, table)
	if err != nil {
		return err
	}
	filename = params.SaveDir + "\\" + table.DataFileName
	utility.CreateFileIfNotExists(filename)
	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	for _, tuple := range tuples {
		recording.WriteTableRecord(file, tuple, -1)
	}
	fmt.Printf("Додано %d записів\n", len(tuples))
	return nil
}

func deleteTable(name string) {
	err := SysCatalog.DeleteTableByName(name)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("таблицю %s видалено успішно, нова кількість таблиць: %d", name, len(SysCatalog.Tables))
	}
}

func setWorkdir(pathToDir string) {
	params.WorkDir = pathToDir
}

func setSaveDir(pathToDir string) {
	params.SaveDir = pathToDir
}

func printTables() {
	fmt.Printf("Кількість таблиць: %d\n", len(SysCatalog.Tables))
	for _, table := range SysCatalog.Tables {
		fmt.Print(table.ToString())
		fmt.Print("\n\n\n")
	}
}

func printRecords(tableName string) {
	records := recording.GetRecords(tableName)

	for _, record := range records {
		for _, fieldValue := range record {
			switch fieldValue.ValueType {
			case types.Int_t:
				fmt.Printf("id: %d, value: %d;\n", fieldValue.ID, fieldValue.Value)
			case types.Real_t:
				fmt.Printf("id: %d, value: %d;\n", fieldValue.ID, fieldValue.Value)
			case types.String_t:
			case types.Char_t:
			case types.Color_t:
				fmt.Printf("id: %d, value: %s;\n", fieldValue.ID, fieldValue.Value)
			}
		}
	}
}
