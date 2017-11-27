package api

type Version struct {
	Authors []string `json:"authors"`
	Version string   `json:"version"`
}

func NewVersion(version string, authors []string) *Version {
	return &Version{authors, version}
}
