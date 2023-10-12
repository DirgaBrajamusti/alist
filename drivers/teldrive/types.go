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

// For Uploading
type PartFile struct {
	Name       string `json:"name"`
	PartId     int    `json:"partId"`
	PartNo     int    `json:"partNo"`
	TotalParts int    `json:"totalParts"`
	Size       int64  `json:"size"`
}
type UploadFile struct {
	Parts []PartFile `json:"parts,omitempty"`
}

type FilePart struct {
	ID int `json:"id"`
}

type FileUploadRequest struct {
	Name     string     `json:"name"`
	MimeType string     `json:"mimeType"`
	Type     string     `json:"type"`
	Parts    []FilePart `json:"parts"`
	Size     int        `json:"size"`
	Path     string     `json:"path"`
}
