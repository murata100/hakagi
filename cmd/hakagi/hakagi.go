package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/syucream/hakagi/src/database"
	"github.com/syucream/hakagi/src/formatter"
	"github.com/syucream/hakagi/src/guess"
)

var ruleToGuesser = map[string]guess.GuessOption{
	"primarykey":     guess.GuessByPrimaryKey(),
	"tableandcolumn": guess.GuessByTableAndColumn(),
}

func main() {
	dbUser := flag.String("dbuser", "", "database user")
	dbPass := flag.String("dbpass", "", "database password")
	dbHost := flag.String("dbhost", "localhost", "database host")
	dbPort := flag.Int("dbport", 3306, "database port")

	targets := flag.String("targets", "", "analysing target databases(comma-separated)")
	rules := flag.String("rules", "primarykey,tableandcolumn", "analysing rules(comma-separated)")

	format := flag.String("format", "sql", "output format(sql / xml)")
	qBase := flag.String("qbase", "indexes", "query base(indexes / columns)")
	output := flag.String("output", "", "output to file")
	cTypes := flag.String("ctypes", "", "compatible types(idtype:columntype)")

	flag.Parse()

	db, err := database.ConnectDatabase(*dbUser, *dbPass, *dbHost, *dbPort)
	if err != nil {
		log.Fatalf("Failed to connect database : %v", err)
	}

	isIndexesQueryBase := true
	if *qBase == "indexes" {
		isIndexesQueryBase = true
	} else if *qBase == "columns" {
		isIndexesQueryBase = false
	} else {
		log.Fatalf("Unknown query base")
	}
	targetSlice := strings.Split(*targets, ",")
	schemas, err := database.FetchSchemas(db, targetSlice, isIndexesQueryBase)
	if err != nil {
		log.Fatalf("Failed to fetch schemas : %v", err)
	}
	primaryKeys, err := database.FetchPrimaryKeys(db, targetSlice)
	if err != nil {
		log.Fatalf("Failed to fetch primary keys : %v", err)
	}

	var guessers []guess.GuessOption
	for _, rule := range strings.Split(*rules, ",") {
		if guesser, ok := ruleToGuesser[rule]; ok {
			guessers = append(guessers, guesser)
		}
	}
	compatibleTypes := strings.Split(*cTypes, ":")
	constraints := guess.GuessConstraints(schemas, primaryKeys, compatibleTypes, guessers...)

	outputText := ""
	if *format == "sql" {
		outputText = formatter.FormatSql(constraints)
	} else if *format == "xml" {
		outputText = formatter.FormatXML(constraints)
	} else {
		log.Fatalf("Unknown output format")
	}

	if *output == "" {
		fmt.Println(outputText)
	} else {
		err := ioutil.WriteFile(*output, ([]byte)(outputText), 0666)
		if err != nil {
			log.Fatalf("Failed to write : %v", err)
		}
	}
}
