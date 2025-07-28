package gsheet

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/tpryan/work/artifact"
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

// New returns a new GSheet object to act as a datasource
func New(svc sheets.Service, sheetID string) GSheet {
	g := GSheet{svc: &svc, id: sheetID}
	return g

}

// SheetID returns the sheet ID for the individual sheet
// (think tabs at the bottom) for a given spreadsheet
func (g *GSheet) SheetID(name string) (int64, error) {

	ranges := []string{name}

	resp, err := g.svc.Spreadsheets.Get(g.id).Ranges(ranges...).Context(context.Background()).Do()
	if err != nil {
		if strings.Contains(err.Error(), "Unable to parse range") {
			return 0, errGSheetDoesNotExist
		}
		return 0, err
	}

	return resp.Sheets[0].Properties.SheetId, nil
}

var errGSheetDoesNotExist = fmt.Errorf("sheets: input sheet does not exist")
var errGSheetAlreadyExists = fmt.Errorf("sheets: input sheet already exists")

// Clear removes all content from an input sheet name
func (g *GSheet) Clear(name string) error {

	clearRange := fmt.Sprintf("%s!A:Z", name)

	req := &sheets.ClearValuesRequest{}

	if _, err := g.svc.Spreadsheets.Values.Clear(g.id, clearRange, req).Do(); err != nil {
		if strings.Contains(err.Error(), "Unable to parse range") {
			return errGSheetDoesNotExist
		}
		return err
	}

	return nil
}

// Add creates a new sheet in the spreadsheet with the input name
func (g *GSheet) Add(name string) error {

	rbb := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				AddSheet: &sheets.AddSheetRequest{
					Properties: &sheets.SheetProperties{
						Title: name,
					},
				},
			}},
	}

	if _, err := g.svc.Spreadsheets.BatchUpdate(g.id, rbb).Do(); err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("A sheet with the name \"%s\" already exists", name)) {
			return errGSheetAlreadyExists
		}
		return fmt.Errorf("sheets: failed to create a new sheet %s", err)
	}

	return nil
}

// Delete removes a sheet in the spreadsheet with the input name
func (g *GSheet) Delete(name string) error {

	id, err := g.SheetID(name)
	if err != nil {
		return err
	}

	rbb := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				DeleteSheet: &sheets.DeleteSheetRequest{
					SheetId: id,
				},
			}},
	}

	if _, err := g.svc.Spreadsheets.BatchUpdate(g.id, rbb).Do(); err != nil {
		return fmt.Errorf("sheets: failed to delete sheet %s", err)
	}
	return nil
}

// FormatSheet creates a set of batch requests that will format a sheet for
// displaying artifacts
func (g *GSheet) FormatSheet(id int64) []*sheets.Request {

	batchreq := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
					Fields: "gridProperties.frozenRowCount",
					Properties: &sheets.SheetProperties{
						SheetId: id,
						GridProperties: &sheets.GridProperties{
							FrozenRowCount: 1,
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat",
					Range: &sheets.GridRange{
						SheetId:          id,
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndRowIndex:      1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColorStyle",
					Range: &sheets.GridRange{
						SheetId:          id,
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndRowIndex:      1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColorStyle: &sheets.ColorStyle{
								RgbColor: &sheets.Color{
									Red:   .85,
									Blue:  1.0,
									Green: .85,
								}},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.numberFormat",
					Range: &sheets.GridRange{
						SheetId:          id,
						StartColumnIndex: 5,
						EndColumnIndex:   6,
						StartRowIndex:    1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							NumberFormat: &sheets.NumberFormat{
								Type:    "DATE",
								Pattern: "mm/dd/yyyy",
							},
						},
					},
				},
			},
			{
				AutoResizeDimensions: &sheets.AutoResizeDimensionsRequest{
					Dimensions: &sheets.DimensionRange{
						SheetId:    id,
						Dimension:  "COLUMNS",
						StartIndex: 0,
						EndIndex:   6,
					},
				},
			},
		},
	}

	return batchreq.Requests
}

// FormatRows generates batch requests to format individual rows of a row
// where the rows consist of set of Artifacts
func (g *GSheet) FormatRows(id int64, a artifact.Artifacts) []*sheets.Request {

	batchreq := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{},
	}

	for i, art := range a {

		if art.Subproject == "" {
			req := &sheets.Request{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColorStyle",
					Range: &sheets.GridRange{
						SheetId:          id,
						StartColumnIndex: 0,
						StartRowIndex:    int64(i) + 1,
						EndRowIndex:      int64(i) + 2,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColorStyle: &sheets.ColorStyle{
								RgbColor: &sheets.Color{
									Red:   1.0,
									Blue:  .98,
									Green: .98,
								}},
						},
					},
				},
			}

			batchreq.Requests = append(batchreq.Requests, req)
		}

		if art.Type == "" {
			req := &sheets.Request{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColorStyle",
					Range: &sheets.GridRange{
						SheetId:          id,
						StartColumnIndex: 0,
						StartRowIndex:    int64(i) + 1,
						EndRowIndex:      int64(i) + 2,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColorStyle: &sheets.ColorStyle{
								RgbColor: &sheets.Color{
									Red:   1.0,
									Blue:  .50,
									Green: .95,
								}},
						},
					},
				},
			}

			batchreq.Requests = append(batchreq.Requests, req)
		}

		if art.Project == "" {
			req := &sheets.Request{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColorStyle",
					Range: &sheets.GridRange{
						SheetId:          id,
						StartColumnIndex: 0,
						StartRowIndex:    int64(i) + 1,
						EndRowIndex:      int64(i) + 2,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColorStyle: &sheets.ColorStyle{
								RgbColor: &sheets.Color{
									Red:   1.0,
									Blue:  .90,
									Green: .90,
								}},
						},
					},
				},
			}

			batchreq.Requests = append(batchreq.Requests, req)
		}

	}

	return batchreq.Requests
}

// ToSheet sends an interface to the named Sheet
func (g *GSheet) ToSheet(name string, i Interfacer) error {

	id, err := g.SheetID(name)
	if err == nil {
		if err := g.Clear(name); err != nil {
			return fmt.Errorf("sheets: failed to clear sheet %s", err)
		}
	}

	if err != nil {
		if err != errGSheetDoesNotExist {
			return fmt.Errorf("sheets: failed to clear sheet %s", err)
		}
		if err := g.Add(name); err != nil {
			return err
		}
		id, err = g.SheetID(name)
		if err != nil {
			return err
		}
	}

	if err := g.UpdateData(name, i); err != nil {
		return fmt.Errorf("sheets: failed to insert into sheet %s", err)
	}

	batchreq := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat",
					Range: &sheets.GridRange{
						SheetId: id,
					},
				},
			},
		},
	}

	batchreq.Requests = append(batchreq.Requests, g.FormatSheet(id)...)
	batchreq.Requests = append(batchreq.Requests, g.FormatRows(id, i.(artifact.Artifacts))...)

	if _, err := g.svc.Spreadsheets.BatchUpdate(g.id, batchreq).Do(); err != nil {
		return fmt.Errorf("sheets: failed to apply formatting %s", err)
	}

	return nil
}

// UpdateData inserts a given set of interfacer data into the spreadsheet in
// sheet name
func (g *GSheet) UpdateData(name string, i Interfacer) error {

	var vr sheets.ValueRange
	vr.Values = i.ToInterfaces()

	r := fmt.Sprintf("%s!A%d:Z100000", name, 1)

	if _, err := g.svc.Spreadsheets.Values.Update(g.id, r, &vr).ValueInputOption("USER_ENTERED").Do(); err != nil {
		if strings.Contains(err.Error(), "Unable to parse range") {
			return errGSheetDoesNotExist
		}

		return fmt.Errorf("sheets: failed to insert data into sheet: %s", err)
	}
	return nil
}

// Artifacts returns a given sheet as Artifacts
func (g *GSheet) Artifacts(name string) (artifact.Artifacts, error) {
	as := artifact.Artifacts{}
	ranges := []string{name}

	resp, err := g.svc.Spreadsheets.Get(g.id).Ranges(ranges...).IncludeGridData(true).Do()
	if err != nil {
		if strings.Contains(err.Error(), "Unable to parse range") {
			return nil, errGSheetDoesNotExist
		}

		return nil, fmt.Errorf("sheets: couldn't read from spreadsheet: %w", err)

	}
	for i, row := range resp.Sheets[0].Data[0].RowData {

		if i == 0 {
			continue
		}
		if len(row.Values) < 7 {
			continue
		}

		as = append(as, newArtifact(row))

	}

	return as, nil
}

func newArtifact(row *sheets.RowData) artifact.Artifact {
	a := artifact.Artifact{}
	a.Type = strings.ReplaceAll(extractString(*row.Values[0]), "\n", "")
	a.Project = extractString(*row.Values[1])
	a.Subproject = extractString(*row.Values[2])
	a.Title = strings.ReplaceAll(extractString(*row.Values[3]), "\n", "")
	a.Role = extractString(*row.Values[4])
	a.ShippedDate = extractTime(*row.Values[5])
	a.Link = extractString(*row.Values[6])
	return a
}

func extractString(val sheets.CellData) string {
	if val.EffectiveValue != nil && val.EffectiveValue.StringValue != nil {
		return strings.TrimSpace(*val.EffectiveValue.StringValue)
	}
	return ""
}

func extractTime(val sheets.CellData) time.Time {
	sqlformat := "2006-01-02 15:04:05.999999-07"
	otherformat := "01/02/2006"
	if val.EffectiveValue != nil {
		if val.EffectiveValue.NumberValue != nil {
			d := int(*val.EffectiveValue.NumberValue) - 2
			epochDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
			result := epochDate.AddDate(0, 0, d)
			return result
		}

		if val.EffectiveValue.StringValue != nil {
			result, err := time.Parse(sqlformat, *val.EffectiveValue.StringValue)
			if err == nil {
				return result
			}

			result2, err := time.Parse(otherformat, *val.EffectiveValue.StringValue)
			if err == nil {
				return result2
			}
		}
		tmp, err := val.EffectiveValue.MarshalJSON()
		if err != nil {
			log.Errorf("could not get json for: (%v) :%s", val.EffectiveValue, err)
			return time.Time{}
		}
		log.Warnf("Was asked to translate a time and wasn't able to: %v", string(tmp))
	}

	return time.Time{}
}
