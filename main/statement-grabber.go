package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ggilmore/politifact-request-maker/config"
	"github.com/ggilmore/politifact-request-maker/scraper"
	"github.com/ggilmore/politifact-request-maker/workers"

	"github.com/zabawaba99/firego"
)

var peopleJSON, _ = filepath.Abs("../politifact-request-maker/resources/people.json")
var subjectsJSON, _ = filepath.Abs("../politifact-request-maker/resources/subjects.json")

const USAGE = "[METHOD (\"person\" or \"subject\")] [NAME_RESOURCE] [MAX_ITEMS] [OUTPUT_DIR] \n OR \n \"serve\" [CONFIG_FILE_PATH]"

func handleStatements(method scraper.StatementMethod, args []string) {
	if len(args) != 5 {
		log.Fatalf("Need 4 args \n%v", USAGE)

	}

	name := args[2]
	n, err := strconv.Atoi(args[3])

	if err != nil {

		log.Fatalf("3rd arg not an int %v, \n %v", err, USAGE)
	}

	dir := args[4]

	var bytes int64
	var path string

	switch method {
	case scraper.ByPerson:
		cleanName := scraper.NameSlugFromFile(name, peopleJSON)
		statements := scraper.SortBySubject(scraper.StatementRequest(scraper.ByPerson, cleanName, n))
		path = dir + cleanName + "-statements.json"
		bytes = scraper.WriteSortedStatementFile(statements, path)

	case scraper.BySubject:
		cleanName := scraper.NameSlugFromFile(name, subjectsJSON)
		statements := scraper.StatementRequest(scraper.BySubject, cleanName, n)
		path = dir + cleanName + "-statements.json"
		bytes = scraper.WriteStatementFile(statements, path)

	default:
		fmt.Println(USAGE)
		os.Exit(1)
	}

	log.Print("Wrote %v bytes to: %v path .", strconv.FormatInt(bytes, 10), path)
	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Invalid arg list: %v", USAGE)
	}

	switch os.Args[1] {
	case "person":
		handleStatements(scraper.ByPerson, os.Args)
	case "subject":
		handleStatements(scraper.BySubject, os.Args)
	case "serve":
		fName := os.Args[2]
		config, err := config.LoadConfig(fName)
		if err != nil {
			log.Fatalf("%v \n %v", USAGE, err)
		}

		inbound := make(chan []scraper.Statement)
		outbound := make(chan scraper.Statement)

		fb := firego.New(config.Firebase.Root, nil)

		go workers.DiffStatements(inbound, outbound)

		go workers.MakePolitifactRequests(config.Politifact.RequestRate, config.Politifact.RequestSize, inbound)

		workers.SendSubjectStatements(config.Firebase.MaxConcurrentReqs, config.Firebase.PeopleChildName, config.Firebase.SubjectsChildName, outbound, fb)

	default:
		log.Fatal(USAGE)

	}

}
