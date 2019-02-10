package schema

// Backend represents an elasticsearch backend
type Backend struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Simple   bool   `yaml:"simple"`
}
