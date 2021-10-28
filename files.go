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
	`{short_hostname}`,
	`{short_hostname}-{instance}`,
	`{short_hostname}-{deployment}`,
	`{short_hostname}-{deployment}-{instance}`,
	`{full_hostname}`,
	`{full_hostname}-{instance}`,
	`{full_hostname}-{deployment}`,
	`{full_hostname}-{deployment}-{instance}`,
	`local`,
	`local-{instance}`,
	`local-{deployment}`,
	`local-{deployment}-{instance}`,
}

const envFileName = "env"

var defaultProviders = map[string]provider.Provider{
	".yaml": &provider.YamlProvider{},
	".yml":  &provider.YamlProvider{},
}

func getExpectedBasename(tmpl string, env environment) (ret string) {
	ret = strings.Replace(tmpl, "{deployment}", env.deployment, 1)

	return ret
}
