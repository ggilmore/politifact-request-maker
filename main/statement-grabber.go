package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ggilmore/wevote-scraper/scraper"
)

var peopleJSON, _ = filepath.Abs("../wevote-scraper/resources/people.json")
var subjectsJSON, _ = filepath.Abs("../wevote-scraper/resources/subjects.json")

const USAGE = "[METHOD (PERSON OR SUBJECT)] [NAME_RESOURCE] [MAX_ITEMS] [OUTPUT_DIR]"

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Need 3 args.")
		fmt.Println(USAGE)
		os.Exit(1)
	}

	method := os.Args[1]
	name := os.Args[2]
	n, err := strconv.Atoi(os.Args[3])

	if err != nil {
		fmt.Println("3rd arg not an int", err)
		fmt.Println(USAGE)
		os.Exit(1)
	}

	dir := os.Args[4]

	var bytes int64
	var path string

	switch method {
	case "person":
		cleanName := scraper.NameSlugFromFile(name, peopleJSON)
		statements := scraper.SortBySubject(scraper.StatementRequest(scraper.ByPerson, cleanName, n))
		path = dir + cleanName + "-statements.json"
		bytes = scraper.WriteSortedStatementFile(statements, path)

	case "subject":
		cleanName := scraper.NameSlugFromFile(name, subjectsJSON)
		statements := scraper.StatementRequest(scraper.BySubject, cleanName, n)
		path = dir + cleanName + "-statements.json"
		bytes = scraper.WriteStatementFile(statements, path)

	default:
		fmt.Println(USAGE)
		os.Exit(1)
	}

	fmt.Println("Wrote " + strconv.FormatInt(bytes, 10) + " bytes to: " + path + ".")
	os.Exit(0)
}
