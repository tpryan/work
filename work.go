// Package work defines code for using a gsheet as a datasource and destination
// for work related artifacts
package work

import (
	"fmt"
	"os"
	"time"

	"github.com/tpryan/work/artifact"

	"gopkg.in/yaml.v2"
)

// Config is the collection of settings that will direct artifact collection
type Config struct {
	SpreadSheetID string               `yaml:"spread_sheet_id,omitempty"`
	GithubUser    string               `yaml:"github_user,omitempty"`
	Destinations  Destinations         `yaml:"destinations,omitempty"`
	Sources       []string             `yaml:"sources,omitempty"`
	Classifiers   artifact.Classifiers `yaml:"classifiers,omitempty"`
	QueryDrive    bool                 `yaml:"query_drive,omitempty"`
}

// NewConfig returna a config from a given path
func NewConfig(path string) (*Config, error) {
	config := Config{}

	dat, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("couldn't read the config file: %s", err)
	}

	if err := yaml.Unmarshal(dat, &config); err != nil {
		return nil, fmt.Errorf("couldn't parse the config file: %s", err)
	}

	return &config, nil

}

// Destination is a place to write a report based on the criteria
type Destination struct {
	Sheet    string   `yaml:"sheet,omitempty"`
	Sort     string   `yaml:"sort,omitempty"`
	Summary  bool     `yaml:"summary,omitempty"`
	Criteria Criteria `yaml:"criteria,omitempty"`
}

// Destinations is a collection of destination items
type Destinations []Destination

// Criteria are the filters to match a Destination
type Criteria struct {
	Start   time.Time `yaml:"start,omitempty"`
	End     time.Time `yaml:"end,omitempty"`
	Project string    `yaml:"project,omitempty"`
}
