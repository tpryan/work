package work

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

const dateformat = "01-02-2006"

// Artifact represents a work product.
type Artifact struct {
	Title       string    `yaml:"title,omitempty"`
	Link        string    `yaml:"link,omitempty"`
	Type        string    `yaml:"type,omitempty"`
	Project     string    `yaml:"project,omitempty"`
	Subproject  string    `yaml:"subproject,omitempty"`
	Role        string    `yaml:"role,omitempty"`
	ShippedDate time.Time `yaml:"shipped_date,omitempty"`
	Extra       string    `yaml:"extra,omitempty"`
}

// Copy returns an exact duplicate of an artifact
func (a Artifact) Copy() Artifact {
	return Artifact{
		Type:        a.Type,
		Project:     a.Project,
		Subproject:  a.Subproject,
		Title:       a.Title,
		Role:        a.Role,
		ShippedDate: a.ShippedDate,
		Link:        a.Link,
	}
}

// String returns a string representation of an artifact
func (a Artifact) String() string {
	return fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s,%s",
		a.Type,
		a.Project,
		a.Subproject,
		a.Title,
		a.Role,
		a.ShippedDate.Format(dateformat),
		a.Link,
	)
}

// Hyperlink formats artifact to be a Google Sheet hyperlink
func (a Artifact) Hyperlink() string {
	u := a.Link
	title := a.Link

	if strings.Contains(a.Link, "critique.corp.google.com") {
		title = strings.ReplaceAll(a.Link, "//critique.corp.google.com/", "//")
	}

	if strings.Contains(a.Link, "buganizer.corp.google.com") {
		title = strings.ReplaceAll(a.Link, "//buganizer.corp.google.com/issues", "//b")
	}
	if strings.Contains(a.Link, "b.corp.google.com") {
		title = strings.ReplaceAll(a.Link, "//b.corp.google.com/issues", "//b")
	}

	if strings.Contains(a.Link, "docs.google.com") {
		uo, err := url.Parse(a.Link)
		if err == nil {
			title = fmt.Sprintf("%s://%s%s", uo.Scheme, uo.Host, uo.Path)
			title = strings.ReplaceAll(title, "/edit", "")
			title = strings.ReplaceAll(title, "/view", "")
		}

	}

	return fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", u, title)
}

type Artifacts []Artifact
