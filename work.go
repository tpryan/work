// Package work defines code for using a Google Sheet as a datasource and destination
// for work related artifacts.
package work

import (
	"fmt"
	"os"
	"time"

	"github.com/tpryan/work/internal/artifact"

	"gopkg.in/yaml.v2"
)

// Config represents the collection of settings that will direct artifact collection.
type Config struct {
	SpreadSheetID string               `yaml:"spread_sheet_id,omitempty"`
	GithubUser    string               `yaml:"github_user,omitempty"`
	Destinations  Destinations         `yaml:"destinations,omitempty"`
	Sources       []string             `yaml:"sources,omitempty"`
	Classifiers   artifact.Classifiers `yaml:"classifiers,omitempty"`
	QueryDrive    bool                 `yaml:"query_drive,omitempty"`
}

// NewConfig returns a Config from a given path.
func NewConfig(path string) (*Config, error) {
	config := Config{}

	dat, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("couldn't read the config file: %w", err)
	}

	if err := yaml.Unmarshal(dat, &config); err != nil {
		return nil, fmt.Errorf("couldn't parse the config file: %w", err)
	}

	return &config, nil

}

// Destination represents a place to write a report based on the criteria.
type Destination struct {
	Sheet    string   `yaml:"sheet,omitempty"`
	Sort     string   `yaml:"sort,omitempty"`
	Summary  bool     `yaml:"summary,omitempty"`
	Criteria Criteria `yaml:"criteria,omitempty"`
}

// Destinations represents a collection of Destination items.
type Destinations []Destination

// Criteria represents the filters used to match a Destination.
type Criteria struct {
	Start   Date   `yaml:"start,omitempty"`
	End     Date   `yaml:"end,omitempty"`
	Project string `yaml:"project,omitempty"`
}

// Date is a custom type for handling YYYY-MM-DD dates in YAML.
type Date struct {
	time.Time
}

// UnmarshalYAML handles parsing of YYYY-MM-DD date strings.
func (d *Date) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	if s == "" {
		return nil
	}

	formats := []string{
		"2006-01-02",
		"2006-1-2",
		time.RFC3339,
	}

	var t time.Time
	var err error
	for _, f := range formats {
		t, err = time.Parse(f, s)
		if err == nil {
			d.Time = t
			return nil
		}
	}

	return fmt.Errorf("could not parse date %q: %w", s, err)
}
