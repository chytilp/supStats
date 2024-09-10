package convertor

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/chytilp/supStats/common"
	"github.com/chytilp/supStats/request"
	"github.com/go-playground/assert/v2"
)

func createConfig() *common.Config {
	cwd, _ := filepath.Abs(".")
	config := common.Config{
		BackendUrl:    "",
		FrontendUrl:   "",
		MobileUrl:     "",
		DataFolder:    path.Join(cwd, "data", "new"),
		OldDataFolder: path.Join(cwd, "data", "old"),
	}
	return &config
}

func removeTestFile(testfilePath string) error {
	err := os.Remove(testfilePath)
	if err != nil {
		return err
	}
	return nil
}

func TestTransformFile(t *testing.T) {
	convertor := Convertor{
		config: createConfig(),
	}

	newFilePath, err := convertor.TransformFile("data_2024_02_01.json")
	if err != nil {
		t.Fatal(err)
	}
	defer removeTestFile(newFilePath)
	data, err := request.UnmarshalFromFile[request.OutputData](newFilePath)
	if err != nil {
		t.Fatal(err)
	}
	fe := data.FindItem("front-end/kóder")
	assert.Equal(t, fe.OfferCount, 96)
	assert.Equal(t, len(fe.Children), 3)
	assert.Equal(t, data.FindItem("javascript").OfferCount, 60)
	assert.Equal(t, data.FindItem("typescript").OfferCount, 45)
	be := data.FindItem("back-end")
	assert.Equal(t, be.OfferCount, 164)
	assert.Equal(t, len(be.Children), 17)
	assert.Equal(t, data.FindItem("c#").OfferCount, 12)
	assert.Equal(t, data.FindItem("python").OfferCount, 42)
	mb := data.FindItem("mobilní vývoj")
	assert.Equal(t, mb.OfferCount, 34)
	assert.Equal(t, len(mb.Children), 5)
	assert.Equal(t, data.FindItem("android").OfferCount, 16)
	assert.Equal(t, data.FindItem("flutter").OfferCount, 3)
	assert.Equal(t, data.FindItem("dummy").OfferCount, 0)
	date := time.Date(2024, time.Month(2), 1, 0, 0, 0, 0, time.Local)
	assert.Equal(t, data.DownloadedAt, date)
}

func TestTransformFiles(t *testing.T) {
	convertor := Convertor{
		config: createConfig(),
	}
	filename := "data_2024_02_01.json"
	badFile := "blbost.json"
	inputFiles := []string{filename, badFile}
	transformResult := convertor.TransformFiles(inputFiles)
	defer removeTestFile(transformResult.OutputFiles[filename])
	cwd, _ := filepath.Abs(".")
	outputFilePath := path.Join(cwd, "data", "new", "2024-02", filename)
	assert.Equal(t, len(transformResult.OutputFiles), 1)
	assert.Equal(t, transformResult.OutputFiles[filename], outputFilePath)
	assert.Equal(t, len(transformResult.Errors), 1)
	assert.Equal(t, transformResult.Errors[badFile], errors.New("wrong filename format: blbost.json"))
	assert.Equal(t, len(transformResult.InputFiles), 2)
}
