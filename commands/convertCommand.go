package commands

import (
	"fmt"
	"os"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"

	"github.com/chytilp/supStats/common"
	"github.com/chytilp/supStats/convertor"
)

type ConvertCommand struct {
	config    *common.Config
	inputDir  string
	outputDir string
}

func NewConvertCommand(config *common.Config, inputDir string, outputDir string) ConvertCommand {
	return ConvertCommand{config: config, inputDir: inputDir, outputDir: outputDir}
}

func (c *ConvertCommand) Run() ([]string, error) {
	inputFiles, err := c.getDataFilesOfDir(c.inputDir)
	if err != nil {
		return nil, err
	}
	outputFiles, err := c.getDataFilesOfDir(c.outputDir)
	if err != nil {
		return nil, err
	}
	filesToConvert := c.removeAlreadyExisting(*inputFiles, *outputFiles)
	conv := convertor.NewConvertor(c.config)
	result := conv.TransformFiles(filesToConvert)
	if len(result.Errors) > 0 {
		for file, err := range result.Errors {
			fmt.Printf("Converting file: %s ends with err: %v\n", file, err)
		}
	}
	converted := make([]string, 0, len(result.OutputFiles))
	for file := range result.OutputFiles {
		converted = append(converted, file)
	}
	return converted, nil
}

func (c *ConvertCommand) getDataFilesOfDir(dirPath string) (*[]string, error) {
	folder, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	files, err := folder.Readdir(0)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, 0, len(files))
	for _, v := range files {
		if !v.IsDir() && common.IsCorrectFileFormat(v.Name()) {
			fileNames = append(fileNames, v.Name())
		}
	}
	return &fileNames, nil
}

func (c *ConvertCommand) removeAlreadyExisting(inpFiles []string, outFiles []string) []string {
	filtered := make([]string, 0, len(inpFiles))
	existing := interSection(inpFiles, outFiles)
	for _, f := range inpFiles {
		if !slices.Contains(existing, f) {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

func interSection[T constraints.Ordered](pS ...[]T) []T {
	hash := make(map[T]*int) // value, counter
	result := make([]T, 0)
	for _, slice := range pS {
		duplicationHash := make(map[T]bool) // duplication checking for individual slice
		for _, value := range slice {
			if _, isDup := duplicationHash[value]; !isDup { // is not duplicated in slice
				if counter := hash[value]; counter != nil { // is found in hash counter map
					if *counter++; *counter >= len(pS) { // is found in every slice
						result = append(result, value)
					}
				} else { // not found in hash counter map
					i := 1
					hash[value] = &i
				}
				duplicationHash[value] = true
			}
		}
	}
	return result
}
