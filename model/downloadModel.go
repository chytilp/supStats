package model

import (
	"fmt"
	"time"
)

type FileContentItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type FileContent struct {
	Categories   []FileContentItem `json:"categories"`
	Technologies []FileContentItem `json:"technologies"`
	DownloadedAt time.Time         `json:"downloaded"`
}

func (f *FileContent) DateInString() string {
	date := f.DownloadedAt
	return fmt.Sprintf("%04d-%02d-%02d", date.Year(), date.Month(), date.Day())
}
