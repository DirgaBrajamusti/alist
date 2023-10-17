package ddrivev3

type Item struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Size     int    `json:"size"`
	IsFolder bool   `json:"isFolder"`
	Modified string `json:"modified"`
}

type DirectoryResponse struct {
	Directory []Item `json:"directory"`
	Files     []Item `json:"files"`
}
