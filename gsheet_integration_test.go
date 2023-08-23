package work

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/sheets/v4"
)

var gsheetTestID = os.Getenv("WORK_gsheetTestID")
var gsheetTestIDNoPerms = os.Getenv("WORK_gsheetTestIDNoPerms")
var credsTestPath = "config/test-creds.json"

func getTestCreds() (*sheets.Service, error) {
	ctx := context.Background()

	f, err := os.Open(credsTestPath)
	if err != nil {
		return nil, err
	}

	config, err := NewConfig(ctx, f, []string{"https://www.googleapis.com/auth/spreadsheets"})
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

	sheetsSVC, err := getTestCreds()
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
	sheetsSVC, err := getTestCreds()
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
	sheetsSVC, err := getTestCreds()
	if err != nil {
		t.Fatalf(fmt.Sprintf("Unable to retrieve Sheets client: %v", err))
	}

	tests := map[string]struct {
		id    string
		in    string
		err   error
		clean bool
	}{
		"withError": {
			id:  gsheetTestID,
			in:  "TestAddError",
			err: errGSheetAlreadyExists,
		},
		"success": {
			id:    gsheetTestID,
			in:    "NewSheetThatShouldBeDeleted",
			clean: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			gsheet := NewGSheet(*sheetsSVC, tc.id)

			err := gsheet.Add(tc.in)
			assert.Equal(t, tc.err, err)

			if tc.clean {
				err = gsheet.Delete(tc.in)
				t.Logf("cleanup: delete %s error: %s", tc.in, err)
			}

		})
	}
}

func TestGsheetDelete(t *testing.T) {
	sheetsSVC, err := getTestCreds()
	if err != nil {
		t.Fatalf(fmt.Sprintf("Unable to retrieve Sheets client: %v", err))
	}

	tests := map[string]struct {
		id     string
		in     string
		err    error
		create bool
	}{
		"withError": {
			id:  gsheetTestID,
			in:  "NewSheetThatShouldBeDeletedAndDoesNotExist",
			err: errGSheetDoesNotExist,
		},
		"success": {
			id:     gsheetTestID,
			in:     "NewSheetThatShouldBeDeleted",
			create: true,
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
			assert.Equal(t, tc.err, err)

		})
	}
}
