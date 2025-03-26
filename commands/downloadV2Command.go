package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/chytilp/supStats/common"
)

type ResponseModel struct {
	TotalItems int `json:"totalItems"`
}

type OutputModel struct {
	Parent bool
	Slug   string
	Count  int
}

type DownloadV2Command struct {
	config          *common.Config
	categorySlugs   []string
	technologySlugs []string
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

func (d *DownloadV2Command) Run() (*string, error) {
	outputItems := []OutputModel{}
	for _, slug := range d.categorySlugs {
		output, err := d.processOneItem(slug, true)
		if err != nil {
			fmt.Printf("Error during process slug: %s\n", slug)
			return nil, err
		}
		outputItems = append(outputItems, *output)
		fmt.Printf("Category: %s done\n", slug)
	}
	for _, slug := range d.technologySlugs {
		output, err := d.processOneItem(slug, false)
		if err != nil {
			fmt.Printf("Error during process slug: %s\n", slug)
			return nil, err
		}
		outputItems = append(outputItems, *output)
		fmt.Printf("Technology: %s done\n", slug)
	}
	filePath, err := d.marshalToFile(outputItems)
	if err != nil {
		fmt.Println("err in MarshalToFile")
		return nil, err
	}
	return filePath, nil
}

func (d *DownloadV2Command) marshalToFile(items []OutputModel) (*string, error) {
	categoryItems := []FileContentItem{}
	technologyItems := []FileContentItem{}
	for _, item := range items {
		if item.Parent {
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
	filePath := common.GetWholePath(fileData.DownloadedAt)
	absPath := path.Join(d.config.DataFolder, filePath)
	err = os.WriteFile(absPath, content, 0644)
	if err != nil {
		return nil, err
	}
	return &absPath, nil
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

func (d *DownloadV2Command) processOneItem(slug string, category bool) (*OutputModel, error) {
	var url string
	if category {
		url = d.config.CategoryBaseUrl + slug
	} else {
		url = d.config.TechnologyBaseUrl + slug
	}
	data, err := d.sendRequest(url)
	if err != nil {
		return nil, fmt.Errorf("error during send request: %s", err)
	}
	result, err := d.parseResponse(data)
	if err != nil {
		fmt.Printf("Error during parsing request data: %s\n", err)
		return nil, fmt.Errorf("error during parsing request data: %s", err)
	}
	var respModel ResponseModel = *result
	return &OutputModel{
		Parent: category,
		Slug:   slug,
		Count:  respModel.TotalItems,
	}, nil
}
