package provider

type Provider interface {
	Parse(in []byte) (map[string]interface{}, error)
}
