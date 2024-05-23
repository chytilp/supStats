package stats

import (
	"os"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"

	"github.com/chytilp/supStats/common"
)

func createFilesMap() map[string][]string {
	return map[string][]string{
		"2024-01": {"data_2024_01_05.json", "data_2024_01_15.json", "neco.json"},
		"2024-02": {"data_2024_02_10.json", "data_2024_02_20.json", "neco.json"},
		"2024-03": {"data_2024_03_15.json", "data_2024_03_25.json", "neco.json"},
		"2024-04": {"data_2024_04_01.json", "data_2024_04_11.json", "neco.json"},
		"2024-05": {"data_2024_05_09.json", "data_2024_04_19.json", "neco.json"},
	}
}

func createFiles(rootFolder string, files map[string][]string) error {
	var path string
	var err error
	var f *os.File
	for folder, files := range files {
		path, err = os.MkdirTemp(rootFolder, folder)
		if err != nil {
			return err
		}
		for _, file := range files {
			f, err = os.CreateTemp(path, file)
			f.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createTableCreate(dateFrom time.Time, dateTo time.Time, config *common.Config) TableCreate {
	return TableCreate{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Tech:     Frontend,
		Config:   config,
	}
}

func createDate(parts DatePart) time.Time {
	return time.Date(parts.year, time.Month(parts.month), parts.day, 0, 0, 0, 0, time.Local)
}

type DatePart struct {
	year  int
	month int
	day   int
}

var testCases = []struct {
	fromDate      DatePart
	toDate        DatePart
	expectedFiles []string
}{
	{DatePart{2024, 2, 1}, DatePart{2024, 2, 29}, []string{"data_2024_02_10.json",
		"data_2024_02_20.json"}},
}

func TestTableCreateFileList(t *testing.T) {
	tmpDir := t.TempDir()
	files := createFilesMap()
	err := createFiles(tmpDir, files)
	if err != nil {
		t.Fatal(err)
	}
	config := common.Config{
		DataFolder: tmpDir,
	}
	var tableCreate TableCreate
	for _, test := range testCases {
		tableCreate = createTableCreate(createDate(test.fromDate), createDate(test.toDate), &config)
		tmp, err := tableCreate.fileList()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *tmp, test.expectedFiles)
	}
}
