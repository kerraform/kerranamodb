package provider

type Package struct {
	OS            string       `json:"os"`
	Arch          string       `json:"arch"`
	Filename      string       `json:"filename"`
	DownloadURL   string       `json:"download_url"`
	SHASumsURL    string       `json:"shasums_url"`
	SHASumsSigURL string       `json:"shasums_signature_url"`
	SHASum        string       `json:"shasum"`
	SigningKeys   *SigningKeys `json:"signing_keys"`
}

// Inspired by Terraform client
// https://github.com/hashicorp/terraform/blob/5da30c2b65265c9c7ac7580a1295e87715ebb568/internal/getproviders/registry_client.go#L214
type SigningKeys struct {
	GPGPublicKeys []GPGPublicKey `json:"gpg_public_keys"`
}

type GPGPublicKey struct {
	KeyID      string `json:"key_id"`
	ASCIIArmor string `json:"ascii_armor"`
}
