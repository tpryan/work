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

// Artifacts is a collection of Artifact items
type Artifacts []Artifact

// ToInterfaces converts artifacts to the slice of slice of interfaces format
// that gsheet requires for data input
func (a Artifacts) ToInterfaces() [][]interface{} {
	var result [][]interface{}

	header := []interface{}{"Type", "Project", "Subproject", "Title", "Role", "Shipped Date", "Link"}
	result = append(result, header)

	for _, v := range a {
		myval := []interface{}{v.Type, v.Project, v.Subproject, v.Title, v.Role, v.ShippedDate.Format("01/02/2006"), v.Hyperlink()}
		result = append(result, myval)
	}

	return result
}

// Massage runs through all of the options in a queue to prune an otherwise
// alter the list of artifacts
func (a *Artifacts) Massage(opts ...Option) *Artifacts {
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Option is function that alters a list of Artifacts
type Option = func(a *Artifacts)

// After returns artifacts from after a particular shippedDate
func After(t time.Time) Option {
	return func(a *Artifacts) {
		result := Artifacts{}

		for _, art := range *a {
			if art.ShippedDate.After(t) {
				result = append(result, art)
			}
		}

		*a = result
	}
}

// Before returns artifacts from before a particular shippedDate
func Before(t time.Time) Option {
	return func(a *Artifacts) {
		result := Artifacts{}

		for _, art := range *a {
			if art.ShippedDate.Before(t) {
				result = append(result, art)
			}
		}

		*a = result
	}
}

// Between returns artifacts from before and after particular shippedDates
func Between(start, end time.Time) Option {
	return func(a *Artifacts) {
		result := Artifacts{}
		if start.IsZero() || end.IsZero() {
			return
		}

		for _, art := range *a {
			if art.ShippedDate.After(start) && art.ShippedDate.Before(end) {
				result = append(result, art)
			}
		}

		*a = result
	}

}

// ProjectFilter returns only the input projects
func ProjectFilter(project string) Option {
	return func(a *Artifacts) {
		if project == "" {
			return
		}
		result := Artifacts{}

		for _, art := range *a {

			if strings.ToLower(art.Project) == strings.ToLower(project) {
				result = append(result, art)
			}

		}

		*a = result
	}
}

// Unique removes repeated artifacts based on links
func Unique() Option {
	return func(a *Artifacts) {
		result := Artifacts{}
		uniquer := map[string]Artifact{}

		for _, art := range *a {
			alr, ok := uniquer[art.Link]
			if !ok {
				uniquer[art.Link] = art
				continue
			}
			if alr.Role != "assignee" && art.Role == "assignee" {
				uniquer[art.Link] = art
				continue
			}

		}

		for _, v := range uniquer {
			result = append(result, v)
		}

		*a = result
	}
}

// ExcludeTitle removes articles that have the input string in the title
func ExcludeTitle(s string) Option {
	return func(a *Artifacts) {
		result := Artifacts{}

		for _, art := range *a {
			if !strings.Contains(art.Title, s) {
				result = append(result, art)
			}
		}

		*a = result
	}
}
