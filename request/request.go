package request

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

func SendRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("create request error: %s\n", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
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

func ParseResponse(data []byte) (*ResponseData, error) {
	var responseData ResponseData
	err := json.Unmarshal(data, &responseData)
	if err != nil {
		return nil, err
	}
	return &responseData, nil
}

func UnmarshalFromFile[T any](filepath string) (*T, error) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var responseData T
	err = json.Unmarshal(byteValue, &responseData)
	if err != nil {
		return nil, err
	}
	return &responseData, nil
}

func MarshalToFile(data OutputData, config *common.Config, version int) (*string, error) {
	content, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return nil, err
	}
	filePath := common.GetWholePath(data.DownloadedAt, version)
	absPath := path.Join(config.DataFolder, filePath)
	err = os.WriteFile(absPath, content, 0644)
	if err != nil {
		return nil, err
	}
	return &absPath, nil
}
