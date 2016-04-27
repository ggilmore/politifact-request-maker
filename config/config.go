package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//Config - config info for the whole program
type Config struct {
	Politifact Politifact `json:"politifact"`
	Firebase   Firebase   `json:"firebase"`
}

//Politifact - configs related to the Politifact API
type Politifact struct {
	RequestSize int      `json:"request_size"`
	RequestRate int      `json:"request_rate_seconds"`
	RSSNames    []string `json:"rss_names"`
}

//Firebase - configs related to the Firebase API
type Firebase struct {
	Root              string `json:"root"`
	PeopleChildName   string `json:"people_child_name"`
	SubjectsChildName string `json:"subjects_child_name"`
	StoriesChildName  string `json:"stories_child_name"`
	MaxConcurrentReqs int    `json:"max_concurrent_requests"`
}

//LoadConfig - load config file, returns a config struct, or an error about what went wrong
//fName - name of the config file
func LoadConfig(fName string) (Config, error) {
	var config Config

	file, err := ioutil.ReadFile(fName)
	if err != nil {
		return config, fmt.Errorf("Couldn't load file: %v", err)
	}

	err = json.Unmarshal(file, &config)

	if err != nil {
		return config, fmt.Errorf("Couldn't unmarshal config file: %v", err)
	}

	return config, nil
}
