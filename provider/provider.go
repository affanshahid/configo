package provider

// Provider is an interface for parsing a map from binary data
type Provider interface {
	Parse(in []byte) (map[string]interface{}, error)
}
