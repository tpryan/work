package work

import (
	"context"

	"google.golang.org/api/sheets/v4"
)

// Interfacer describes a type that can self convert to Interfaces{}{} for
// adding data to Sheets
type Interfacer interface {
	ToInterfaces() [][]interface{}
}

// GSheet provides read/write access to a Google Sheet. Given the correctly
// initialized service it basically turns a Gsheet into a datasource
type GSheet struct {
	svc *sheets.Service
	id  string
}

// NewGSheet returns a new GSheet object to act as a datasource
func NewGSheet(svc sheets.Service, sheetID string) GSheet {
	g := GSheet{svc: &svc, id: sheetID}
	return g

}

// SheetID returns the sheet ID for the individual sheet
// (think tabs at the bottom) for a given spreadsheet
func (g *GSheet) SheetID(name string) (int64, error) {

	ranges := []string{name}

	resp, err := g.svc.Spreadsheets.Get(g.id).Ranges(ranges...).Context(context.Background()).Do()
	if err != nil {
		return 0, err
	}

	return resp.Sheets[0].Properties.SheetId, nil
}
