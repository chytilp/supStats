package stats

import (
	"os"
	"strings"
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
		"2024-05": {"data_2024_05_09.json", "data_2024_05_19.json", "neco.json"},
	}
}

func createFiles(rootFolder string, files map[string][]string) error {
	var err error
	var f *os.File
	for folder, files := range files {
		err = os.Mkdir(rootFolder+"/"+folder, 0755)
		if err != nil {
			return err
		}
		for _, file := range files {
			f, err = os.Create(rootFolder + "/" + folder + "/" + file)
			f.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createTableCreate[T Number](dateFrom time.Time, dateTo time.Time, config *common.Config) TableCreate[T] {
	return TableCreate[T]{
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

func shortenPath(wholePath string) string {
	parts := strings.Split(wholePath, "/")
	length := len(parts)
	return parts[length-2] + "/" + parts[length-1]
}

var testCases = []struct {
	fromDate      DatePart
	toDate        DatePart
	expectedFiles []string
}{
	{DatePart{2024, 2, 1}, DatePart{2024, 2, 29}, []string{"2024-02/data_2024_02_10.json",
		"2024-02/data_2024_02_20.json"}},
	{DatePart{2024, 2, 1}, DatePart{2024, 3, 31}, []string{"2024-02/data_2024_02_10.json",
		"2024-02/data_2024_02_20.json", "2024-03/data_2024_03_15.json", "2024-03/data_2024_03_25.json"}},
	{DatePart{2024, 3, 1}, DatePart{2024, 4, 30}, []string{"2024-03/data_2024_03_15.json",
		"2024-03/data_2024_03_25.json", "2024-04/data_2024_04_01.json", "2024-04/data_2024_04_11.json"}},
	{DatePart{2024, 2, 11}, DatePart{2024, 3, 16}, []string{"2024-02/data_2024_02_20.json",
		"2024-03/data_2024_03_15.json"}},
	{DatePart{2024, 4, 11}, DatePart{2024, 5, 19}, []string{"2024-04/data_2024_04_11.json",
		"2024-05/data_2024_05_09.json", "2024-05/data_2024_05_19.json"}},
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
	var tableCreate TableCreate[int]
	for _, test := range testCases {
		tableCreate = createTableCreate[int](createDate(test.fromDate), createDate(test.toDate), &config)
		tmp, err := tableCreate.fileList()
		if err != nil {
			t.Fatal(err)
		}
		soubory := *tmp
		for index, f := range test.expectedFiles {
			assert.Equal(t, shortenPath(soubory[index]), f)
		}
	}
}
