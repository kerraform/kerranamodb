package provider

type AvailableVersion struct {
	Version   string                     `json:"version"`
	Protocols []string                   `json:"protocols"`
	Platforms []AvailableVersionPlatform `json:"platforms"`
}

type AvailableVersionPlatform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}
