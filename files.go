package configo

import (
	"strings"

	"github.com/affanshahid/configo/provider"
)

var orderedTemplates = []string{
	`default`,
	`default-{instance}`,
	`{deployment}`,
	`{deployment}-{instance}`,
	`{shortHostname}`,
	`{shortHostname}-{instance}`,
	`{shortHostname}-{deployment}`,
	`{shortHostname}-{deployment}-{instance}`,
	`{fullHostname}`,
	`{fullHostname}-{instance}`,
	`{fullHostname}-{deployment}`,
	`{fullHostname}-{deployment}-{instance}`,
	`local`,
	`local-{instance}`,
	`local-{deployment}`,
	`local-{deployment}-{instance}`,
}

const envFileName = "env"

var defaultProviders = map[string]provider.Provider{
	".yaml":  &provider.YamlProvider{},
	".yml":   &provider.YamlProvider{},
	".json":  &provider.JsonProvider{},
	".json5": &provider.Json5Provider{},
	".hjson": &provider.HjsonProvider{},
	".toml":  &provider.TomlProvider{},
}

func getExpectedBasename(tmpl string, env environment) (ret string) {
	ret = strings.Replace(tmpl, "{deployment}", env.deployment, 1)

	return ret
}
