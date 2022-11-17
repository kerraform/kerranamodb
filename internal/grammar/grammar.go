package grammar

const (
	// Inspired by Semantic Versioning 2.0.0
	// Ref: https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
	Version = `(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)(?:-(?:(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?:[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`
)
