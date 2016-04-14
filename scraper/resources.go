package scraper

import (
	"encoding/gob"
	"os"
)

func deepcopy(dst, src interface{}) error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(w)
	err = enc.Encode(src)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(r)
	return dec.Decode(dst)
}

type Subject struct {
	Subject     string `json:"subject"`
	SubjectSlug string `json:"subject_slug"`
}

type Person struct {
	Party          Party  `json:"party"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	NameSlug       string `json:"name_slug"`
	CanonicalPhoto string `json:"canonical_photo"`
}

type Party struct {
	Party     string `json:"party"`
	PartySlug string `json:"party_slug"`
}

const PersonEndpoint string = "http://www.politifact.com/api/people/all/json/"

const SubjectEndpoint string = "http://www.politifact.com/api/subjects/all/json/"

type Ruling struct {
	RulingSlug       string `json:"ruling_slug"`
	Ruling           string `json:"ruling"`
	CanonicalGraphic string `json:"canonical_ruling_graphic"`
}

type Art struct {
	Caption        string `json:"caption"`
	CanonicalPhoto string `json:"canonical_photo"`
	Youtube        string `json:"youtube"`
	YoutubeID      string `json:"youtubeID"`
	Title          string `json:"title"`
}
type Story struct {
	UpdatedDate     string `json:"updated_date"`
	Art             Art    `json:"art"`
	Headline        string `json:"headline"`
	PublicationDate string `json:"publication_date"`
	StoryURL        string `json:"story_url"`
	Blurb           string `json:"blurb"`
}

type StatementType struct {
	StatementType string `json:"statement_type"`
}

type Statement struct {
	StatementURL     string        `json:"statement_url"`
	Target           []Person      `json:"target"`
	StatementDate    string        `json:"statement_date"`
	StatementContext string        `json:"statement_context"`
	Speaker          Person        `json:"speaker"`
	RulingHeadline   string        `json:"ruling_headline"`
	Statement        string        `json:"statement"`
	Ruling           Ruling        `json:"ruling"`
	RulingLinkTest   string        `json:"ruling_link_test"`
	RulingDate       string        `json:"ruling_date"`
	StatementType    StatementType `json:"statement_type"`
	Subject          []Subject     `json:"subject"`
}

type StatementMethod int

const (
	ByPerson StatementMethod = iota
	BySubject
)
