package work

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/sheets/v4"
)

var gsheetTestID = os.Getenv("WORK_gsheetTestID")
var gsheetTestIDNoPerms = os.Getenv("WORK_gsheetTestIDNoPerms")
var credsTestPath = "config/test-creds.json"

func getTestSheetsSvc() (*sheets.Service, error) {
	ctx := context.Background()

	f, err := os.Open(credsTestPath)
	if err != nil {
		return nil, err
	}

	config, err := NewClientOption(ctx, f, []string{"https://www.googleapis.com/auth/spreadsheets"})
	if err != nil {
		return nil, err
	}

	sheetsSVC, err := sheets.NewService(ctx, config)
	if err != nil {
		return nil, err
	}

	return sheetsSVC, nil

}

func TestGsheetSheetID(t *testing.T) {

	sheetsSVC, err := getTestSheetsSvc()
	if err != nil {
		t.Fatalf(fmt.Sprintf("Unable to retrieve Sheets client: %v", err))
	}

	tests := map[string]struct {
		id     string
		in     string
		want   int64
		errStr string
	}{
		"withError": {
			id:     gsheetTestIDNoPerms,
			in:     "Manual",
			want:   0,
			errStr: "Error 403: The caller does not have permission",
		},
		"success": {
			id:   gsheetTestID,
			in:   "Manual",
			want: 1766204602,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			gsheet := NewGSheet(*sheetsSVC, tc.id)

			got, err := gsheet.SheetID(tc.in)
			if tc.errStr == "" && err != nil {
				t.Logf("sheet name: %s", tc.in)
				t.Logf("spreadsheet id: %s", tc.id)
				t.Fatalf("got an error when expected none: %s", err)
			}
			if tc.errStr != "" && strings.Contains(err.Error(), tc.errStr) {
				t.Skip()
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGsheetClear(t *testing.T) {
	sheetsSVC, err := getTestSheetsSvc()
	if err != nil {
		t.Fatalf(fmt.Sprintf("Unable to retrieve Sheets client: %v", err))
	}

	tests := map[string]struct {
		id  string
		in  string
		err error
	}{
		"withError": {
			id:  gsheetTestID,
			in:  "TestClearNotExists",
			err: errGSheetDoesNotExist,
		},
		"success": {
			id: gsheetTestID,
			in: "TestClear",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			gsheet := NewGSheet(*sheetsSVC, tc.id)

			err := gsheet.Clear(tc.in)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestGsheetAdd(t *testing.T) {
	sheetsSVC, err := getTestSheetsSvc()
	if err != nil {
		t.Fatalf(fmt.Sprintf("Unable to retrieve Sheets client: %v", err))
	}

	tests := map[string]struct {
		id     string
		in     string
		errStr string
		clean  bool
	}{
		"withError": {
			id:     gsheetTestID,
			in:     "TestAddError",
			errStr: errGSheetAlreadyExists.Error(),
		},
		"success": {
			id:    gsheetTestID,
			in:    "NewSheetThatShouldBeDeleted",
			clean: true,
		},
		"protected": {
			id:     gsheetTestID,
			in:     "TestProtected",
			errStr: errGSheetAlreadyExists.Error(),
		},
		"toolong": {
			id:     gsheetTestID,
			in:     "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789",
			errStr: "The sheet name cannot be greater than 100 characters",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			gsheet := NewGSheet(*sheetsSVC, tc.id)

			err := gsheet.Add(tc.in)
			if tc.errStr != "" && err != nil {
				assert.Contains(t, err.Error(), tc.errStr)
			}

			if tc.clean {
				err = gsheet.Delete(tc.in)
				t.Logf("cleanup: delete %s error: %s", tc.in, err)
			}

		})
	}
}

func TestGsheetDelete(t *testing.T) {
	sheetsSVC, err := getTestSheetsSvc()
	if err != nil {
		t.Fatalf(fmt.Sprintf("Unable to retrieve Sheets client: %v", err))
	}

	tests := map[string]struct {
		id     string
		in     string
		errStr string
		create bool
	}{
		"withError": {
			id:     gsheetTestID,
			in:     "NewSheetThatShouldBeDeletedAndDoesNotExist",
			errStr: errGSheetDoesNotExist.Error(),
		},
		"success": {
			id:     gsheetTestID,
			in:     "NewSheetThatShouldBeDeleted",
			create: true,
		},
		"protected": {
			id:     gsheetTestID,
			in:     "TestProtected",
			errStr: "You are trying to edit a protected cell or object",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			gsheet := NewGSheet(*sheetsSVC, tc.id)

			if tc.create {
				err := gsheet.Add(tc.in)
				t.Logf("create: %s error: %s", tc.in, err)
			}

			err = gsheet.Delete(tc.in)

			if tc.errStr == "" {
				assert.Nil(t, err)
				t.Skip()
			}

			assert.ErrorContains(t, err, tc.errStr)

		})
	}
}

func TestGSheetArtifacts(t *testing.T) {

	sheetsSVC, err := getTestSheetsSvc()
	if err != nil {
		t.Fatalf(fmt.Sprintf("Unable to retrieve Sheets client: %v", err))
	}

	tests := map[string]struct {
		id     string
		in     string
		want   Artifacts
		errStr string
	}{
		"basic": {
			id: gsheetTestID,
			in: "TestArtifacts",
			want: Artifacts{
				Artifact{
					Type:        "Website",
					Project:     "DeployStack",
					Subproject:  "Core Platform",
					Role:        "Primary",
					ShippedDate: time.Date(2021, 12, 2, 0, 0, 0, 0, time.UTC),
					Link:        "https://appinabox.dev/",
				},
			},
		},
		"BadName": {
			id:     gsheetTestID,
			in:     "TestArtifactsShouldBeBadNAme",
			want:   Artifacts{},
			errStr: "input sheet does not exist",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			gsheet := NewGSheet(*sheetsSVC, tc.id)

			got, err := gsheet.Artifacts(tc.in)
			if tc.errStr == "" && err != nil {
				t.Fatalf("got an error when expected none: %s", err)
			}

			if tc.errStr != "" {
				assert.ErrorContains(t, err, tc.errStr)
				t.Skip()
			}

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGSheetUpdateData(t *testing.T) {

	sheetsSVC, err := getTestSheetsSvc()
	if err != nil {
		t.Fatalf(fmt.Sprintf("Unable to retrieve Sheets client: %v", err))
	}

	tests := map[string]struct {
		id     string
		name   string
		in     Artifacts
		errStr string
	}{
		"basic": {
			id:   gsheetTestID,
			name: "TestUpdateData",
			in: Artifacts{
				Artifact{
					Title:       time.Now().String(),
					Type:        "test",
					Project:     "Test",
					Subproject:  "Test",
					Role:        "Test",
					ShippedDate: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC),
					Link:        "https://example.com/",
				},
			},
		},
		"basicError": {
			id:     gsheetTestID,
			name:   "TestUpdateDataDoesNotExist",
			in:     nil,
			errStr: errGSheetDoesNotExist.Error(),
		},

		"protected": {
			id:   gsheetTestID,
			name: "TestProtected",
			in: Artifacts{
				Artifact{
					Type:        "Website",
					Project:     "DeployStack",
					Subproject:  "Core Platform",
					Role:        "Primary",
					ShippedDate: time.Date(2021, 12, 2, 0, 0, 0, 0, time.UTC),
					Link:        "https://appinabox.dev/",
				},
			},
			errStr: "You are trying to edit a protected cell or object",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			gsheet := NewGSheet(*sheetsSVC, tc.id)

			err := gsheet.UpdateData(tc.name, tc.in)
			if tc.errStr == "" && err != nil {
				t.Fatalf("got an error when expected none: %s", err)
			}
			if tc.errStr != "" && !strings.Contains(err.Error(), tc.errStr) {
				t.Fatalf("got: '%s' expected : '%s'", err, tc.errStr)
			}

			got, err := gsheet.Artifacts(tc.name)

			assert.Equal(t, tc.in, got)
		})
	}
}

func TestGSheetToSheet(t *testing.T) {

	sheetsSVC, err := getTestSheetsSvc()
	if err != nil {
		t.Fatalf(fmt.Sprintf("Unable to retrieve Sheets client: %v", err))
	}

	tests := map[string]struct {
		id     string
		name   string
		in     Artifacts
		errStr string
		clean  bool
	}{
		"basic": {
			id:   gsheetTestID,
			name: "TestToSheet",
			in: Artifacts{
				Artifact{
					Title:       time.Now().String(),
					Type:        "test",
					Project:     "Test",
					Subproject:  "Test",
					Role:        "Test",
					ShippedDate: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC),
					Link:        "https://example.com/",
				},
			},
		},
		"basicWithCreate": {
			id:   gsheetTestID,
			name: "TestToSheetDoesNotExist",
			in: Artifacts{
				Artifact{
					Title:       time.Now().String(),
					Type:        "test",
					Project:     "Test",
					Subproject:  "Test",
					Role:        "Test",
					ShippedDate: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC),
					Link:        "https://example.com/",
				},
			},
			clean: true,
		},
		"empty": {
			id:    gsheetTestID,
			name:  "TestToSheetDoesNotExistEmpty",
			in:    Artifacts{},
			clean: true,
		},
		"nameTooLong": {
			id:     gsheetTestID,
			name:   "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789",
			in:     nil,
			errStr: "The sheet name cannot be greater than 100 characters",
		},
		"protected": {
			id:   gsheetTestID,
			name: "TestProtected",
			in: Artifacts{
				Artifact{
					Type:        "Website",
					Project:     "DeployStack",
					Subproject:  "Core Platform",
					Role:        "Primary",
					ShippedDate: time.Date(2021, 12, 2, 0, 0, 0, 0, time.UTC),
					Link:        "https://appinabox.dev/",
				},
			},
			errStr: "You are trying to edit a protected cell or object",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			gsheet := NewGSheet(*sheetsSVC, tc.id)

			err := gsheet.ToSheet(tc.name, tc.in)
			if tc.errStr == "" && err != nil {
				t.Fatalf("got an error when expected none: %s", err)
			}
			if tc.errStr != "" && !strings.Contains(err.Error(), tc.errStr) {
				t.Fatalf("got: '%s' expected : '%s'", err, tc.errStr)
			}

			got, err := gsheet.Artifacts(tc.name)

			assert.Equal(t, tc.in, got)

			if tc.clean {
				err = gsheet.Delete(tc.name)
				t.Logf("cleanup: delete %s error: %s", tc.in, err)
			}

		})
	}
}
