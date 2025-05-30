package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/chytilp/supStats/commands"
	"github.com/chytilp/supStats/common"
	"github.com/chytilp/supStats/persistence"
	"github.com/chytilp/supStats/stats"
)

func getDb(config *common.Config) (*sql.DB, error) {
	dbFile := config.DbPath
	err := persistence.CreateSupDatabase(dbFile)
	if err != nil {
		fmt.Printf("CreateDatabase err: %v\n", err)
		return nil, err
	}
	db, err := persistence.GetDatabase(dbFile)
	if err != nil {
		fmt.Printf("GetDatabase err: %v\n", err)
		return nil, err
	}
	return db, nil
}

func main() {
	fmt.Println("Start app")
	config := common.GetConfig()

	tableCmd := flag.NewFlagSet("table", flag.ExitOnError)
	var tableType string
	tableCmd.StringVar(&tableType, "type", "", "type (fe, be, mb)")
	tableCmd.StringVar(&tableType, "t", "", "type (fe, be, mb)")
	var fromTo string
	tableCmd.StringVar(&fromTo, "fromTo", "", "fromTo")
	tableCmd.StringVar(&fromTo, "f", "", "fromTo")
	var columns int
	tableCmd.IntVar(&columns, "columns", 0, "columns")
	tableCmd.IntVar(&columns, "c", 0, "columns")
	var aggColumns bool
	tableCmd.BoolVar(&aggColumns, "aggColumns", false, "aggColumns")
	tableCmd.BoolVar(&aggColumns, "a", false, "aggColumns")

	relTableCmd := flag.NewFlagSet("reltable", flag.ExitOnError)
	var relTableType string
	relTableCmd.StringVar(&relTableType, "type", "", "type (fe, be, mb)")
	relTableCmd.StringVar(&relTableType, "t", "", "type (fe, be, mb)")
	var relFromTo string
	relTableCmd.StringVar(&relFromTo, "fromTo", "", "fromTo")
	relTableCmd.StringVar(&relFromTo, "f", "", "fromTo")
	var relColumns int
	relTableCmd.IntVar(&relColumns, "columns", 0, "columns")
	relTableCmd.IntVar(&relColumns, "c", 0, "columns")

	convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)
	var inputDir string
	convertCmd.StringVar(&inputDir, "inputDir", "", "inputDir")
	convertCmd.StringVar(&inputDir, "i", "", "inputDir")

	importCmd := flag.NewFlagSet("import", flag.ExitOnError)
	var importInputDir string
	importCmd.StringVar(&importInputDir, "inputDir", "", "inputDir")
	var version int
	importCmd.IntVar(&version, "version", 0, "version")

	if len(os.Args) < 2 {
		fmt.Println("expected 'download', 'table', 'relTable' or 'convert' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "download":
		downloadCmd := commands.NewDownloadV2Command(config)
		filePath, err := downloadCmd.Run()
		if err != nil {
			fmt.Println("err in DownloadV3Command")
			log.Fatalln(err.Error())
		}
		fmt.Printf("File %s was saved\n", *filePath)

	case "table":
		tableCmd.Parse(os.Args[2:])
		technology := stats.TechnologyFromString(tableType)
		tableCommand := commands.NewTableCommand(config, technology, fromTo, columns, aggColumns)
		err := tableCommand.Run()
		if err != nil {
			fmt.Println("err in DownloadCommand")
			log.Fatalln(err.Error())
		}

	case "reltable":
		relTableCmd.Parse(os.Args[2:])
		fmt.Println("reltable subcommand")
		fmt.Printf("type: %s, fromTo: %s, columns: %d\n", relTableType, relFromTo, relColumns)

	case "convert":
		convertCmd.Parse(os.Args[2:])
		convertCommand := commands.NewConvertCommand(config, inputDir)
		converted, err := convertCommand.Run()
		if err != nil {
			fmt.Println("err in ConvertCommand")
			log.Fatalln(err.Error())
		}
		for _, convertedFile := range converted {
			fmt.Printf("file: %s was successfully converted.\n", convertedFile)
		}
	case "import":
		importCmd.Parse(os.Args[2:])
		db, err := getDb(config)
		if err != nil {
			fmt.Println("err in create and get database")
			log.Fatalln(err.Error())
		}
		importCommand := commands.NewImportCommand(db, importInputDir, version)
		results := importCommand.Run()
		fmt.Printf("result: %v\n", results)

	default:
		fmt.Println("expected 'download', 'table', 'relTable', 'import' or 'convert' subcommands")
		os.Exit(1)
	}
}
