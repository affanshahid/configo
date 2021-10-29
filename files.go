package configo

import (
	"strings"
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

var defaultProviders = map[string]Provider{
	".yaml":  yamlProvider,
	".yml":   yamlProvider,
	".json":  jsonProvider,
	".json5": json5Provider,
	".hjson": hjsonProvider,
	".toml":  tomlProvider,
}

func getExpectedBasename(tmpl string, env environment) (ret string) {
	ret = strings.Replace(tmpl, "{deployment}", env.deployment, 1)
	ret = strings.Replace(ret, "{instance}", env.instance, 1)
	ret = strings.Replace(ret, "{shortHostname}", env.shortHostname, 1)
	ret = strings.Replace(ret, "{fullHostname}", env.fullHostname, 1)

	return ret
}
