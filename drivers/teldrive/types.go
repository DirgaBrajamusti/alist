package teldrive

import (
	"time"
)

type Session struct {
	UserName string `json:"userName"`
	Hash     string `json:"hash"`
}

type File struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	MimeType  string    `json:"mimeType"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	Starred   bool      `json:"starred"`
	ParentID  string    `json:"parentId"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type FileList struct {
	Results []File `json:"results"`
}
