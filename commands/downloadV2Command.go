package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/chytilp/supStats/common"
)

type ResponseModel struct {
	TotalItems int `json:"totalItems"`
}

type DownloadV2Command struct {
	config          *common.Config
	categorySlugs   []string
	technologySlugs []string
}

type OutputModel struct {
	IsCategory bool
	Slug       string
	Count      int
	Error      *error
}

type FileContentItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type FileContent struct {
	Categories   []FileContentItem `json:"categories"`
	Technologies []FileContentItem `json:"technologies"`
	DownloadedAt time.Time         `json:"downloaded"`
}

func (o *OutputModel) IsValid() bool {
	return o.Error == nil
}

type InputModel struct {
	Slug       string
	IsCategory bool
	BaseUrl    string
}

func (i *InputModel) GetUrl() string {
	return i.BaseUrl + i.Slug
}

func NewDownloadV2Command(config *common.Config) DownloadV2Command {
	category := []string{
		"back-end-vyvojar",
		"front-end-vyvojar",
		"full-stack-vyvojar",
		"mobilni-vyvojar",
		"devops-specialista",
	}
	technology := []string{
		"javascript", "html", "css-3", "typescript", "java", "php", "net",
		"python", "c-sharp", "c%2B%2B", "c", "ruby", "go", "scala", "r", "nodejs",
		"rust", "kotlin", "flutter", "react-native", "dart", "react-js",
		"angular", "vue-js", "docker", "kubernetes", "ansible", "terraform",
		"django", "flask",
	}

	return DownloadV2Command{config: config, categorySlugs: category, technologySlugs: technology}
}

func (d *DownloadV2Command) createInputs() []InputModel {
	inputs := make([]InputModel, len(d.categorySlugs)+len(d.technologySlugs))
	var index int = 0
	for _, categorySlug := range d.categorySlugs {
		inputs[index] = InputModel{Slug: categorySlug, IsCategory: true, BaseUrl: d.config.CategoryBaseUrl}
		index += 1
	}
	for _, technologySlug := range d.technologySlugs {
		inputs[index] = InputModel{Slug: technologySlug, IsCategory: false, BaseUrl: d.config.TechnologyBaseUrl}
		index += 1
	}
	return inputs
}

func (d *DownloadV2Command) Run() (*string, error) {
	inputs := d.createInputs()

	inputChan := make(chan InputModel, 4)
	outputChan := make(chan OutputModel)

	results := make([]OutputModel, len(d.categorySlugs)+len(d.technologySlugs))
	var wg sync.WaitGroup
	wg.Add(4)

	go d.process(inputChan, outputChan, "routine-1", &wg)
	go d.process(inputChan, outputChan, "routine-2", &wg)
	go d.process(inputChan, outputChan, "routine-3", &wg)
	go d.process(inputChan, outputChan, "routine-4", &wg)

	go func() {
		wg.Wait()
		close(outputChan)
	}()

	go feed(inputChan, inputs)

	var index int = 0
	for output := range outputChan {
		if output.IsValid() {
			fmt.Printf("Result: %v\n", output)
			results[index] = output
			index += 1
		} else {
			fmt.Printf("Url: %s, Error: %v\n", output.Slug, *output.Error)
		}
	}

	filePath, err := d.marshalToFile(results)
	if err != nil {
		fmt.Println("err in MarshalToFile")
		return nil, err
	}
	return filePath, nil
}

func (d *DownloadV2Command) process(inputChan <-chan InputModel, outputChan chan<- OutputModel, label string, wg *sync.WaitGroup) {
	defer wg.Done()
	var outputErr error
	for input := range inputChan {
		url := input.GetUrl()
		fmt.Printf("Processor: %s works on url: %s\n", label, url)
		data, err := d.sendRequest(url)
		if err != nil {
			outputErr = fmt.Errorf("error during send request: %s", err)
			outputChan <- OutputModel{Slug: input.Slug, Error: &outputErr}
			continue
		}
		result, err := d.parseResponse(data)
		if err != nil {
			outputErr = fmt.Errorf("error during parsing request data: %s", err)
			outputChan <- OutputModel{Slug: input.Slug, Error: &outputErr}
			continue
		}
		var respModel ResponseModel = *result
		fmt.Printf("Processor: %s successfully finished url: %s\n", label, url)
		outputChan <- OutputModel{Slug: input.Slug, Count: respModel.TotalItems, IsCategory: input.IsCategory}
	}
	fmt.Printf("I am out of cyclus processor: %s\n", label)
}

func (d *DownloadV2Command) sendRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("create request error: %s\n", err)
		return nil, err
	}
	req.Header.Set("Accept", "application/ld+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DNT", "1")
	req.Header.Set("Host", d.config.HeaderHost)
	req.Header.Set("Origin", d.config.HeaderOrigin)
	req.Header.Set("Priority", "u=4")
	req.Header.Set("Referer", d.config.HeaderReferer)

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return nil, err
	}
	fmt.Printf("Response status code: %d, reason: %s\n", res.StatusCode, res.Status)
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return nil, err
	}

	return resBody, nil
}

func (d *DownloadV2Command) parseResponse(data []byte) (*ResponseModel, error) {
	var responseData ResponseModel
	err := json.Unmarshal(data, &responseData)
	if err != nil {
		return nil, err
	}
	return &responseData, nil
}

func feed(inputChan chan<- InputModel, inputs []InputModel) {
	for _, input := range inputs {
		inputChan <- input
	}
	close(inputChan)
}

func (d *DownloadV2Command) marshalToFile(items []OutputModel) (*string, error) {
	categoryItems := []FileContentItem{}
	technologyItems := []FileContentItem{}
	for _, item := range items {
		if item.IsCategory {
			categoryItems = append(categoryItems, FileContentItem{Name: item.Slug, Count: item.Count})
		} else {
			if item.Slug == "c%2B%2B" {
				item.Slug = "c++"
			}
			technologyItems = append(technologyItems, FileContentItem{Name: item.Slug, Count: item.Count})
		}
	}

	fileData := FileContent{
		Categories:   categoryItems,
		Technologies: technologyItems,
		DownloadedAt: time.Now(),
	}
	content, err := json.MarshalIndent(fileData, "", " ")
	if err != nil {
		return nil, err
	}
	filePath := common.GetWholePath(fileData.DownloadedAt, 25)
	absPath := path.Join(d.config.DataFolder, filePath)
	err = os.WriteFile(absPath, content, 0644)
	if err != nil {
		return nil, err
	}
	return &absPath, nil
}
