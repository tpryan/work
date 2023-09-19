package work

import (
	"bytes"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"text/template"
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

func uniform(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Search looks for an exact match for a given link in a given set of artifacts
// it adjusts for shortcuts for buganizer and critique
func (a Artifacts) Search(link string) Artifact {
	for _, art := range a {

		if strings.Contains(uniform(art.Link), uniform(link)) ||
			strings.Contains(uniform(link), uniform(art.Link)) {
			return art
		}

		if strings.Contains(link, "critique.") ||
			strings.Contains(link, "buganizer.") ||
			strings.Contains(link, "b/") {
			sl := strings.Split(link, "/")
			tmp := sl[len(sl)-1]

			if strings.Contains(link, tmp) {
				return art
			}

		}

	}

	return Artifact{}
}

// Copy duplicates a list of artifacts completely
func (a Artifacts) Copy() Artifacts {
	result := Artifacts{}

	for _, v := range a {
		result = append(result, v.Copy())
	}
	return result
}

// Sort reorders a list of artifacts by ShippedDate
func (a *Artifacts) Sort() {
	sort.Slice((*a), func(i, j int) bool {
		return (*a)[i].ShippedDate.Before((*a)[j].ShippedDate)
	})
}

// SortReport reorders a list of artifacts by Project, Subproject, Type, then ShippedDate
func (a *Artifacts) SortReport() {
	sort.Slice((*a), func(i, j int) bool {
		if (*a)[i].Project == (*a)[j].Project {
			if (*a)[i].Subproject == (*a)[j].Subproject {
				if (*a)[i].Type == (*a)[j].Type {
					return (*a)[i].ShippedDate.Before((*a)[j].ShippedDate)
				}
				return (*a)[i].Type < (*a)[j].Type
			}
			return (*a)[i].Subproject < (*a)[j].Subproject
		}
		return (*a)[i].Project < (*a)[j].Project
	})
}

// Template spits out a list of artifacts as a markdown report
func (a Artifacts) Template(label string) (string, error) {
	a.SortReport()
	a.FillInSubs()

	ds := map[string]map[string]map[string][]string{}

	for _, art := range a {
		if _, ok := ds[art.Project]; !ok {
			ds[art.Project] = map[string]map[string][]string{}
		}
		if _, ok := ds[art.Project][art.Subproject]; !ok {
			ds[art.Project][art.Subproject] = map[string][]string{}
		}

		ds[art.Project][art.Subproject][art.Type] = append(ds[art.Project][art.Subproject][art.Type], art.Link)
	}

	type report struct {
		Title   string
		Details map[string]map[string]map[string][]string
	}

	r := report{label, ds}

	text := `# {{.Title}}
{{range $Project, $Subproject  := .Details}}
## {{$Project}}
{{range $Subproject_text, $Type  := $Subproject}}
### {{$Subproject_text}}
{{range $Type_text, $Link  := $Type}}
#### {{$Type_text}}
{{range $Url  := $Link}}
* {{$Url}}{{end}}
{{end}}
{{end}}
{{end}}

`

	t, err := template.New("md").Parse(text)
	if err != nil {
		return "", fmt.Errorf("could not parse template: %s", err)
	}

	var tpl bytes.Buffer

	if err := t.Execute(&tpl, r); err != nil {
		return "", fmt.Errorf("could not execute template: %s", err)
	}

	return tpl.String(), nil
}

// FillInSubs adds N/A for all empty subprojects, for reporting purposes
func (a *Artifacts) FillInSubs() {
	result := Artifacts{}

	for _, art := range *a {
		if art.Subproject == "" {
			art.Subproject = "N/A"
		}
		result = append(result, art)
	}
	*a = result
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

// Classify analyzes a set of artifacts and fills in Project and Subproject
// based on matches or substrings
func Classify(list Classifiers) Option {
	return func(a *Artifacts) {
		result := Artifacts{}

	ArtLoop:
		for _, art := range *a {

			for _, v := range list.Exclusions {
				if strings.Contains(art.Link, v) {
					continue ArtLoop
				}
			}

			// Find if the item is in the classify list somewhere
			class := list.Search(art.Link)

			// otherwise if it matches overwrite and continue
			if class.Link != "" {
				art.Project = class.Project
				art.Subproject = class.Subproject
				result = append(result, art)
				continue
			}

			art = list.Stamp(art)

			result = append(result, art)

		}
		*a = result
	}

}

// Classifier is a data structure that is used for filling in missing data in
// artifacts
type Classifier struct {
	Project    string              `yaml:"project,omitempty"`
	Subproject string              `yaml:"subproject,omitempty"`
	Links      []string            `yaml:"links,omitempty"`
	Contains   map[string][]string `yaml:"contains,omitempty"`
}

// Classifiers is a collection of Classifer items
type Classifiers struct {
	Lists      []Classifier `yaml:"lists,omitempty"`
	Exclusions []string     `yaml:"exclusions,omitempty"`
	artifacts  Artifacts
}

// Search loons through a list of classifiers and returns a Artifact template
// to use in filling in missing data in the items that match the link
func (c Classifiers) Search(link string) Artifact {
	if c.artifacts == nil {
		result := Artifacts{}

		for _, list := range c.Lists {
			for _, link := range list.Links {
				na := Artifact{}
				na.Project = list.Project
				na.Subproject = list.Subproject
				na.Link = link
				result = append(result, na)
			}
		}
		c.artifacts = result
	}

	return c.artifacts.Search(link)
}

// Stamp alters the input artifact based on substring matching
func (c Classifiers) Stamp(art Artifact) Artifact {

	for _, list := range c.Lists {
		for key, value := range list.Contains {
			if key == "title" {
				for _, v := range value {
					if strings.Contains(uniform(art.Title), uniform(v)) {
						art.Project = list.Project
						art.Subproject = list.Subproject
					}
				}
			}
			if key == "link" {
				for _, v := range value {
					if strings.Contains(uniform(art.Link), uniform(v)) {
						art.Project = list.Project
						art.Subproject = list.Subproject
					}
				}
			}
		}
	}

	return art
}
