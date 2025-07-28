package gsheet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tpryan/work/artifact"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// var gsheetTestID = os.Getenv("WORK_gsheetTestID")
var gsheetTestID = "1T3DDzZCSXp31uG6yY_sc_IRmnfLFxrHIKCbZi6noDRM"
var gsheetTestIDNoPerms = os.Getenv("WORK_gsheetTestIDNoPerms")
var credsTestPath = "../testdata/test-creds.json"

// TODO: Dedupe this function
// NewClientOption returns a clientOption from a given set of credentials.
// Used to initialize Google API clients
func NewClientOption(ctx context.Context, r io.Reader, scopes []string) (option.ClientOption, error) {
	creds := struct {
		ClientEmail  string `json:"client_email"`
		PrivateKey   string `json:"private_key"`
		PrivateKeyID string `json:"private_key_id"`
		TokenURL     string `json:"token_uri"`
	}{}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, fmt.Errorf("could not read credentials: %w", err)
	}

	if err := json.Unmarshal(buf.Bytes(), &creds); err != nil {
		return nil, fmt.Errorf("error unmarshaling credentials file: %w", err)
	}

	conf := &jwt.Config{
		Email:        creds.ClientEmail,
		PrivateKey:   []byte(creds.PrivateKey),
		PrivateKeyID: creds.PrivateKeyID,
		TokenURL:     creds.TokenURL,
		Scopes:       scopes,
	}

	client := option.WithHTTPClient(conf.Client(ctx))

	return client, nil
}

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
		t.Fatalf("Unable to retrieve Sheets client: %v", err)
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
			t.Logf("envsheetid: %s", tc.id)
			gsheet := NewGSheet(*sheetsSVC, tc.id)
			t.Logf("sheet id: %v", tc.id)

			got, err := gsheet.SheetID(tc.in)
			if tc.errStr == "" && err != nil {
				t.Logf("sheet name: %s", tc.in)
				t.Logf("spreadsheet id: %s", tc.id)
				t.Fatalf("got an error when expected none: %s", err)
			}
			if tc.errStr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errStr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGsheetClear(t *testing.T) {
	sheetsSVC, err := getTestSheetsSvc()
	if err != nil {
		t.Fatalf("Unable to retrieve Sheets client: %v", err)
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
		t.Fatalf("Unable to retrieve Sheets client: %v", err)
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
		t.Fatalf("Unable to retrieve Sheets client: %v", err)
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
		t.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	tests := map[string]struct {
		id     string
		in     string
		want   artifact.Artifacts
		errStr string
	}{
		"basic": {
			id: gsheetTestID,
			in: "TestArtifacts",
			want: artifact.Artifacts{
				artifact.Artifact{
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
			want:   artifact.Artifacts{},
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
		t.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	tests := map[string]struct {
		id     string
		name   string
		in     artifact.Artifacts
		errStr string
	}{
		"basic": {
			id:   gsheetTestID,
			name: "TestUpdateData",
			in: artifact.Artifacts{
				artifact.Artifact{
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
			in: artifact.Artifacts{
				artifact.Artifact{
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
			if tc.errStr == "" {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.in, got)
		})
	}
}

func TestGSheetToSheet(t *testing.T) {

	sheetsSVC, err := getTestSheetsSvc()
	if err != nil {
		t.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	tests := map[string]struct {
		id     string
		name   string
		in     artifact.Artifacts
		errStr string
		clean  bool
	}{
		"basic": {
			id:   gsheetTestID,
			name: "TestToSheet",
			in: artifact.Artifacts{
				artifact.Artifact{
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
			in: artifact.Artifacts{
				artifact.Artifact{
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
			in:    artifact.Artifacts{},
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
			in: artifact.Artifacts{
				artifact.Artifact{
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
			if tc.errStr == "" {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.in, got)

			if tc.clean {
				err = gsheet.Delete(tc.name)
				t.Logf("cleanup: delete %s error: %s", tc.in, err)
			}

		})
	}
}
