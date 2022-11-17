package driver

type DriverType string

const (
	DriverTypeLocal DriverType = "local"
	DriverTypeS3    DriverType = "s3"
)

type Driver interface{}

type CreateModuleVersionResult struct {
	Upload string
}

type ProviderVersionMetadata struct {
	KeyID string `json:"key-id"`
}

type CreateProviderVersionResult struct {
	SHASumsUpload    string
	SHASumsSigUpload string
}

type CreateProviderPlatformResult struct {
	ProviderBinaryUploads string
}
