package migrate

// Databases returns the configured logical database names in config order.
//
// It intentionally exposes only names, not configured source strings, database
// URL strings, or resolved secret values.
func (s *Migrator) Databases() []string {
	databases := make([]string, 0, len(s.config.Databases))
	for _, d := range s.config.Databases {
		databases = append(databases, d.Name)
	}

	return databases
}
