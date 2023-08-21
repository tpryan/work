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

const gsheetTestID = "1T3DDzZCSXp31uG6yY_sc_IRmnfLFxrHIKCbZi6noDRM"
const gsheetTestIDNoPerms = "1q3Bqa4BuOMrYfQHRCJohoHVR_ks8sVbteFP1xoJy3sY"
const credsTestPath = "config/test-creds.json"

func TestGsheetSheetID(t *testing.T) {
	ctx := context.Background()

	f, err := os.Open(credsTestPath)
	if err != nil {
		t.Fatalf("couldn't read the credential file: %s", err)
	}

	config, err := NewConfig(ctx, f, []string{"https://www.googleapis.com/auth/spreadsheets"})
	if err != nil {
		t.Fatalf("could not generate a config: %s", err)
	}

	sheetsSVC, err := sheets.NewService(ctx, config)
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
			want:   1766204602,
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
