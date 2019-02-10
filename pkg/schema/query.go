package schema

// Query represents an elasticsearch query. Supported types are the default,
// which is go-templated json, yaml, as well as lua-scripted
type Query struct {
	Name    string `yaml:"name"`
	Backend string `yaml:"backend"`
	Index   string `yaml:"index"`
	Type    string `yaml:"type"`
	Source  string `yaml:"source"`
}
