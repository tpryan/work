package work

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/sheets/v4"
)

func TestGsheetFormatSheet(t *testing.T) {
	tests := map[string]struct {
		in   int64
		want []*sheets.Request
	}{
		"basic": {
			in: 1,
			want: []*sheets.Request{
				{
					UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
						Fields: "gridProperties.frozenRowCount",
						Properties: &sheets.SheetProperties{
							SheetId: 1,
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
							SheetId:          1,
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
							SheetId:          1,
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
							SheetId:          1,
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
							SheetId:    1,
							Dimension:  "COLUMNS",
							StartIndex: 0,
							EndIndex:   6,
						},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tmp := GSheet{}
			got := tmp.FormatSheet(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGsheetFormatRow(t *testing.T) {
	tests := map[string]struct {
		in        int64
		artifacts Artifacts
		want      []*sheets.Request
	}{
		"basic": {
			in: 1,
			artifacts: Artifacts{
				Artifact{
					Project:    "Project",
					Subproject: "Subproject",
					Title:      "Title",
					Type:       "Type",
				},
			},
			want: []*sheets.Request{},
		},
		"NoSubproject": {
			in: 1,
			artifacts: Artifacts{
				Artifact{
					Project: "Project",
					Title:   "Title",
					Type:    "Type",
				},
			},
			want: []*sheets.Request{
				&sheets.Request{
					RepeatCell: &sheets.RepeatCellRequest{
						Fields: "userEnteredFormat.backgroundColorStyle",
						Range: &sheets.GridRange{
							SheetId:          1,
							StartColumnIndex: 0,
							StartRowIndex:    int64(1),
							EndRowIndex:      int64(2),
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
				},
			},
		},
		"NoType": {
			in: 1,
			artifacts: Artifacts{
				Artifact{
					Project:    "Project",
					Subproject: "Subproject",
					Title:      "Title",
				},
			},
			want: []*sheets.Request{
				{
					RepeatCell: &sheets.RepeatCellRequest{
						Fields: "userEnteredFormat.backgroundColorStyle",
						Range: &sheets.GridRange{
							SheetId:          1,
							StartColumnIndex: 0,
							StartRowIndex:    int64(1),
							EndRowIndex:      int64(2),
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
				},
			},
		},
		"NoProject": {
			in: 1,
			artifacts: Artifacts{
				Artifact{
					Subproject: "Subproject",
					Title:      "Title",
					Type:       "Type",
				},
			},
			want: []*sheets.Request{
				{
					RepeatCell: &sheets.RepeatCellRequest{
						Fields: "userEnteredFormat.backgroundColorStyle",
						Range: &sheets.GridRange{
							SheetId:          1,
							StartColumnIndex: 0,
							StartRowIndex:    int64(1),
							EndRowIndex:      int64(2),
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
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tmp := GSheet{}
			got := tmp.FormatRows(tc.in, tc.artifacts)
			assert.Equal(t, tc.want, got)
		})
	}
}
