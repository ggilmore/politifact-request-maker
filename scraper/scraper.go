package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ggilmore/politifact-request-maker/resources"
)

//normalizes all names to the same format
func sanitizeName(s string) string {
	return strings.ToLower(strings.Replace(s, " ", "", -1))
}

//default implementation
func handleError(e error) {
	if e != nil {
		panic(e)
	}
}

//WritePersonFile writes a json file of people
//map of person_slug -> person
func WritePersonFile(ps []Person, fName string) int {

	nameMap := make(map[string]string)

	for _, p := range ps {
		nameMap[(sanitizeName(p.FirstName + p.LastName))] = p.NameSlug
	}

	json, jerr := json.MarshalIndent(nameMap, "", " ")
	handleError(jerr)

	f, createErr := os.Create(fName)
	handleError(createErr)

	defer f.Close()

	bytes, writeErr := f.Write(json)
	handleError(writeErr)

	f.Sync()
	return bytes
}

//WriteSubjectFile writes a json file of subjects
//map of subject_slug -> slug
func WriteSubjectFile(subs []Subject, fName string) int {
	subMap := make(map[string]string)

	for _, s := range subs {
		subMap[sanitizeName(s.Subject)] = s.SubjectSlug
	}

	json, jerr := json.MarshalIndent(subMap, "", " ")
	handleError(jerr)

	f, createErr := os.Create(fName)
	handleError(createErr)

	defer f.Close()

	bytes, writeErr := f.Write(json)
	handleError(writeErr)

	f.Sync()
	return bytes

}

func PersonRequest(endpoint string) []Person {
	resp, requestErr := http.Get(endpoint)

	defer resp.Body.Close()

	handleError(requestErr)

	var r []Person
	decoder := json.NewDecoder(resp.Body)

	jsonErr := decoder.Decode(&r)

	handleError(jsonErr)

	return r
}

func MakeSubjectRequest(endpoint string) []Subject {
	resp, requestErr := http.Get(endpoint)

	defer resp.Body.Close()

	handleError(requestErr)

	var r []Subject
	decoder := json.NewDecoder(resp.Body)

	jsonErr := decoder.Decode(&r)

	handleError(jsonErr)

	return r
}

func statementEndp(met StatementMethod, nslg string, n int) string {
	var method string
	switch met {
	case NoMethod:
		method = ""
		nslg = "/"
	case ByPerson:
		method = "people/"
		nslg = nslg + "/"
	case BySubject:
		method = "subjects/"
		nslg = nslg + "/"
	default:
		panic("Unhandled statement method")
	}

	return StatementEndpoint + method + nslg + "json/?n=" + strconv.Itoa(n)
}

func StatementRequest(met StatementMethod, name string, n int) []Statement {
	resp, requestErr := http.Get(statementEndp(met, name, n))

	defer resp.Body.Close()

	handleError(requestErr)

	var r []Statement
	jsonErr := json.NewDecoder(resp.Body).Decode(&r)

	handleError(jsonErr)

	fmt.Println("done")
	return r
}

func StatementsByDate(n int) []Statement {
	resp, err := http.Get(statementEndp(NoMethod, "", n))
	defer resp.Body.Close()

	handleError(err)
	var r []Statement
	jsonErr := json.NewDecoder(resp.Body).Decode(&r)

	handleError(jsonErr)

	return r

}

func SortBySubject(stmts []Statement) map[string][]Statement {
	groupMap := make(map[string][]Statement)

	for _, stmt := range stmts {
		for _, sub := range stmt.Subject {
			list, _ := groupMap[sub.SubjectSlug]

			list = append(list, stmt)
			groupMap[sub.SubjectSlug] = list

		}
	}

	return groupMap
}

func WriteStatementFile(ss []Statement, fName string) int64 {
	f, createErr := os.Create(fName)
	handleError(createErr)

	defer f.Close()

	jerr := json.NewEncoder(f).Encode(ss)
	handleError(jerr)

	f.Sync()
	stat, err := f.Stat()
	handleError(err)
	return stat.Size()
}

func WriteSortedStatementFile(ss map[string][]Statement, fName string) int64 {
	f, createErr := os.Create(fName)
	handleError(createErr)

	defer f.Close()

	jerr := json.NewEncoder(f).Encode(ss)
	handleError(jerr)

	f.Sync()
	stat, err := f.Stat()
	handleError(err)
	return stat.Size()
}

func NameSlug(name string) string {

	cleanName := sanitizeName(name)

	return resources.PeopleSlugMap[cleanName]

}
