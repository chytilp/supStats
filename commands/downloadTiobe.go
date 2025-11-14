package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chytilp/supStats/common"
	"github.com/chytilp/supStats/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/jfbus/httprs"
)

type DownloadTiobeCommand struct {
	config *common.Config
}

func NewDownloadTiobeCommand(config *common.Config) DownloadTiobeCommand {
	return DownloadTiobeCommand{config: config}
}

func (d *DownloadTiobeCommand) Run() (*string, error) {
	content, err := d.downloadData("https://www.tiobe.com/tiobe-index/")
	defer content.Close()
	if err != nil {
		fmt.Printf("Err at Download data: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Data downloaded.")
	records, err := d.parseTop20Table(content)
	if err != nil {
		fmt.Printf("Err at ParseTop20Table: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Top20 table parsed.")
	fmt.Printf("Records: %d\n", len(records))
	fmt.Printf("First: %s\n", records[0].String())
	fmt.Printf("Middle: %s\n", records[9].String())
	fmt.Printf("Last: %s\n", records[len(records)-1].String())

	content.Seek(0, io.SeekStart)
	other, err := d.parseOtherTable(content)
	if err != nil {
		fmt.Printf("Err at ParseOtherTable: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Other table parsed.")
	fmt.Printf("Records: %d\n", len(other))
	fmt.Printf("First: %s\n", other[0].String())
	fmt.Printf("Middle: %s\n", other[14].String())
	fmt.Printf("Last: %s\n", other[len(other)-1].String())

	content.Seek(0, io.SeekStart)
	monthArrPtr, err := d.parseMonth(content)
	if err != nil {
		fmt.Printf("Err at ParseMonth: %v\n", err)
		os.Exit(1)
	}
	monthArr := *monthArrPtr
	finalDir := fmt.Sprintf("%s-%s", monthArr[0], monthArr[1])
	fmt.Printf("Month parsed: %s\n", finalDir)
	filename := fmt.Sprintf("tiobe_%s_%s.json", monthArr[0], monthArr[1])
	allRecords := append(records, other...)
	dataFile := filepath.Join(d.config.IndexDataFolder, finalDir, filename)
	WriteJsonData(allRecords, dataFile)
	return &dataFile, nil
}

func (d *DownloadTiobeCommand) parseMonth(content io.Reader) (*[2]string, error) {
	doc, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		fmt.Printf("Failed to parse the HTML document: %v\n", err)
		return nil, err
	}
	h1ElText := doc.Find("article.content h1 b").First().Text()
	parts := strings.Split(h1ElText, " ")
	month, err := common.GetMonthNumber(parts[3])
	if err != nil {
		fmt.Printf("Failed to get month number of %s: %v\n", parts[3], err)
		return nil, err
	}
	var monthArr [2]string = [2]string{parts[4], month}
	return &monthArr, nil
}

func (d *DownloadTiobeCommand) parseOtherTable(content io.Reader) ([]model.IndexRow, error) {
	doc, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		fmt.Printf("Failed to parse the HTML document: %v\n", err)
		return nil, err
	}
	otherRows := make([]model.IndexRow, 30)
	var index int = 0
	doc.Find("table#otherPL tr").Each(func(i int, tr *goquery.Selection) {
		if i > 0 {
			var order int
			var lang, sRating string
			var rating float64

			tr.Find("td").Each(func(j int, td *goquery.Selection) {
				switch j {
				case 0:
					order, err = strconv.Atoi(td.Text())
				case 1:
					lang = td.Text()
				case 2:
					sRating = strings.ReplaceAll(td.Text(), "%", "")
					rating, _ = strconv.ParseFloat(sRating, 8)
				}
				if err != nil {
					fmt.Printf("Err parse td value: %v\n", err)
				}
			})
			otherRows[index] = model.NewTiobeOtherRow(lang, order, rating)
			index += 1
		}
	})
	return otherRows, nil
}

func (d *DownloadTiobeCommand) parseTop20Table(content io.Reader) ([]model.IndexRow, error) {
	doc, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		fmt.Printf("Failed to parse the HTML document: %v\n", err)
		return nil, err
	}
	top20Rows := make([]model.IndexRow, 20)
	var index int = 0
	doc.Find("table#top20 tr").Each(func(i int, tr *goquery.Selection) {
		if i > 0 {
			var order, orderPrevYear int
			var lang, sRating, sRatingChange, change string
			var rating, ratingChange float64

			tr.Find("td").Each(func(j int, td *goquery.Selection) {
				switch j {
				case 0:
					order, err = strconv.Atoi(td.Text())
				case 1:
					orderPrevYear, err = strconv.Atoi(td.Text())
				case 2:
					if td.Children().Length() == 0 {
						change = "beze zmeny"
					} else {
						imgEl := td.Children().First()
						if imgEl != nil {
							src, _ := imgEl.Attr("src")
							if strings.HasSuffix(src, "down.png") {
								change = "dolu"
							} else if strings.HasSuffix(src, "up.png") {
								change = "nahoru"
							} else {
								fmt.Printf("Unknown change option: %s\n", src)
								os.Exit(1)
							}
						} else {
							change = "beze zmeny"
						}
					}
				case 4:
					lang = td.Text()
				case 5:
					sRating = strings.ReplaceAll(td.Text(), "%", "")
					rating, err = strconv.ParseFloat(sRating, 8)
				case 6:
					sRatingChange = strings.ReplaceAll(td.Text(), "%", "")
					ratingChange, err = strconv.ParseFloat(sRatingChange, 8)
				}
				if err != nil {
					fmt.Printf("Err parse td value: %v\n", err)
				}
			})
			top20Rows[index] = model.NewTiobeTopRow(lang, order, change, orderPrevYear, rating, ratingChange)
			index += 1
		}
	})
	return top20Rows, nil
}

func (d *DownloadTiobeCommand) downloadData(url string) (io.ReadSeekCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error download data from url: %s, err: %v\n", url, err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("HTTP Error %d: %s\n", resp.StatusCode, resp.Status)
		fmt.Println(msg)
		return nil, fmt.Errorf(msg)
	}
	rs := httprs.NewHttpReadSeeker(resp)
	return rs, nil
}
