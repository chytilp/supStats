package convertor

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/chytilp/supStats/common"
	"github.com/chytilp/supStats/request"
)

type Convertor struct {
	config *common.Config
}

type TransformationResult struct {
	InputFiles  []string
	OutputFiles map[string]string
	Errors      map[string]error
}

func (c *Convertor) TransformFile(oldFilename string) (string, error) {
	oldFilePath := c.config.GetFilePath(true, oldFilename)
	oldFormatData, err := c.readDataFile(oldFilePath)
	if err != nil {
		return "", err
	}
	newFormatData, err := c.transformData(oldFormatData, oldFilename)
	if err != nil {
		return "", err
	}
	newFilePath, err := request.MarshalToFile(*newFormatData, c.config)
	if err != nil {
		return "", err
	}
	return *newFilePath, nil
}

func (c *Convertor) TransformFiles(oldFilenames []string) TransformationResult {
	outputFiles := make(map[string]string, len(oldFilenames))
	errors := make(map[string]error, len(oldFilenames))

	for _, inputFileName := range oldFilenames {
		outputFilePath, err := c.TransformFile(inputFileName)
		if err != nil {
			errors[inputFileName] = err
		} else {
			outputFiles[inputFileName] = outputFilePath
		}
	}

	return TransformationResult{
		InputFiles:  oldFilenames,
		OutputFiles: outputFiles,
		Errors:      errors,
	}
}

func (c *Convertor) readDataFile(oldFilePath string) ([]OldItem, error) {
	data, err := os.ReadFile(oldFilePath)
	if err != nil {
		return nil, err
	}
	var items []OldItem
	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (c *Convertor) createItem(old OldItem, pathName string) request.Item {
	return request.Item{
		Name:          old.Name,
		OfferCount:    old.Count,
		Identifier:    pathName + "/" + strings.ToLower(old.Name),
		UrlIdentifier: "",
	}
}

func (c *Convertor) transformBranch(branch OldItem, pathName string) request.Item {
	children := make([]request.Item, 0, len(*branch.Children))
	for _, child := range *branch.Children {
		children = append(children, c.createItem(child, pathName))
	}

	return request.Item{
		Name:          branch.Name,
		OfferCount:    branch.Count,
		Identifier:    pathName,
		UrlIdentifier: "",
		Children:      children,
	}
}

func (c *Convertor) transformData(data []OldItem, originFileName string) (*request.OutputData, error) {
	date, err := common.GetFileDate(originFileName)
	if err != nil {
		return nil, err
	}
	var feBranch, beBranch, mbBranch request.Item
	for _, item := range data {
		if item.Name == "Vývoj" {
			for _, subItem := range *item.Children {
				if subItem.Name == "Front-End/Kóder" {
					feBranch = c.transformBranch(subItem, "vyvoj/front-end-koder")
				} else if subItem.Name == "Back-End" {
					beBranch = c.transformBranch(subItem, "vyvoj/back-end")
				} else if subItem.Name == "Mobilní vývoj" {
					mbBranch = c.transformBranch(subItem, "vyvoj/mobilni-vyvoj")
				}
			}
		}
	}
	outData := request.OutputData{
		Frontend:     &feBranch,
		Backend:      &beBranch,
		Mobile:       &mbBranch,
		DownloadedAt: *date,
	}
	return &outData, nil
}
