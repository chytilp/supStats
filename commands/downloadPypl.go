package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chytilp/supStats/common"
	"github.com/chytilp/supStats/model"

	"github.com/PuerkitoBio/goquery"
)

func WriteJsonData(data []model.IndexRow, outputFile string) {
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Printf("Err marshal data to json: %v\n", err)
		os.Exit(1)
	}
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Printf("Err write file: %v\n", err)
		os.Exit(1)
	}
}

type DownloadPyplCommand struct {
	sourceFile string
	config     *common.Config
}

func NewDownloadPyplCommand(sourceFile string, config *common.Config) DownloadPyplCommand {
	return DownloadPyplCommand{config: config, sourceFile: sourceFile}
}

func (d *DownloadPyplCommand) Run() (*string, error) {
	content := d.getDataReader(d.sourceFile)
	records, monthArr, err := d.parsePyplTable(content)
	if err != nil {
		fmt.Printf("Err at ParsePyplTable: %v\n", err)
		os.Exit(1)
	}
	month := fmt.Sprintf("%s_%s", monthArr[0], monthArr[1])
	fmt.Printf("Pypl Records: %d\n", len(records))
	fmt.Printf("Month: %s\n", month)
	fmt.Printf("First: %s\n", records[0].String())
	fmt.Printf("Middle: %s\n", records[14].String())
	fmt.Printf("Last: %s\n", records[len(records)-1].String())

	outputFile := fmt.Sprintf("pypl_%s.json", month)
	outputDir := d.config.IndexDataFolder + fmt.Sprintf("%s-%s/", monthArr[0], monthArr[1])
	outputFilePath := filepath.Join(outputDir, outputFile)
	WriteJsonData(records, outputFilePath)

	ok := moveFile(d.sourceFile, d.config.IndexDataFolder)
	fmt.Printf("File: %s was archived: %t.\n", d.sourceFile, ok)
	return &outputFilePath, nil
}

func (d *DownloadPyplCommand) getDataReader(filename string) io.Reader {
	content, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Err at ReadFile: %v\n", err)
		os.Exit(1)
	}
	return content
}

func (d *DownloadPyplCommand) parsePyplTable(content io.Reader) ([]model.IndexRow, [2]string, error) {
	month := [2]string{}
	doc, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		fmt.Printf("Failed to parse the HTML document: %v\n", err)
		return nil, month, err
	}

	pElText := doc.Find("p#countryDate").First().Text()
	elTextParts := strings.Split(pElText, " ")
	monthPart, err := common.GetMonthNumber(elTextParts[1])
	if err != nil {
		return nil, month, err
	}
	month[0] = elTextParts[2]
	month[1] = monthPart
	rows := make([]model.IndexRow, 30)
	var index int = 0
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		if i > 0 {
			var change, lang, sRating, sRatingChange string
			var order int
			var rating, ratingChange float64

			tr.Find("td").Each(func(j int, td *goquery.Selection) {
				switch j {
				case 0:
					order, err = strconv.Atoi(td.Text())
				case 1:
					if td.Children().Length() == 0 {
						change = "beze zmeny"
					} else {
						imgEl := td.Children().First()
						if imgEl != nil {
							src, _ := imgEl.Attr("src")
							if strings.HasSuffix(src, "Down.png") {
								change = "dolu"
							} else if strings.HasSuffix(src, "Up.png") {
								change = "nahoru"
							} else {
								fmt.Printf("Unknown change option: %s\n", src)
								os.Exit(1)
							}
						} else {
							change = "beze zmeny"
						}
					}
				case 2:
					lang = td.Text()
				case 3:
					sRating = strings.ReplaceAll(td.Text(), "%", "")
					sRating = strings.TrimSpace(sRating)
					rating, err = strconv.ParseFloat(sRating, 8)
				case 4:
					sRatingChange = strings.ReplaceAll(td.Text(), "%", "")
					sRatingChange = strings.TrimSpace(sRatingChange)
					ratingChange, err = strconv.ParseFloat(sRatingChange, 8)
					if ratingChange == 0.0 {
						ratingChange = math.Abs(ratingChange)
					}
				}
				if err != nil {
					fmt.Printf("Err parse td value: %v\n", err)
				}
			})
			if index < 30 {
				rows[index] = model.NewPyplRow(lang, order, change, rating, ratingChange)
				index += 1
			}
		}
	})
	return rows, month, nil
}

func moveFile(fileWholePath string, indexDataFolder string) bool {
	filename := filepath.Base(fileWholePath)

	if stat, err := os.Stat(indexDataFolder + "/archiv"); err == nil && stat.IsDir() {
		newWholePath := filepath.Join(indexDataFolder, "archiv", filename)
		err = os.Rename(fileWholePath, newWholePath)
		if err != nil {
			fmt.Printf("File: %s was not archived.Err: %s\n", fileWholePath, err)
			return false
		}
		return true
	}
	fmt.Println("Folder archiv was not found")
	return false
}
