package migrate

// Database for migrate.
type Database struct {
	Name   string `yaml:"name" json:"name" toml:"name"`
	Source string `yaml:"source" json:"source" toml:"source"`
	URL    string `yaml:"url" json:"url" toml:"url"`
}

// Config for migrate.
type Config struct {
	Databases []Database `yaml:"databases" json:"databases" toml:"databases"`
}

// Database by name.
func (c *Config) Database(name string) *Database {
	for _, d := range c.Databases {
		if d.Name == name {
			return &d
		}
	}

	return nil
}
