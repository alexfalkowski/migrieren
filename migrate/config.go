package migrate

// Database for migrate.
type Database struct {
	Name   string `yaml:"name"`
	Source string `yaml:"source"`
	URL    string `yaml:"url"`
}

// Config for migrate.
type Config struct {
	Databases []Database `yaml:"databases"`
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
