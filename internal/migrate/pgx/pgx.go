package pgx

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
)

// ErrInvalidMigrationsTable is returned when pgx migration table options are invalid.
var ErrInvalidMigrationsTable = errors.New("migrate pgx: invalid migrations table")

// Config aliases the golang-migrate pgx driver config.
type Config = pgx.Config

// ParseConfig extracts golang-migrate pgx driver options from u.
//
// Supported pgx5 query parameters:
//   - x-migrations-table: migration table name. When omitted, the upstream
//     driver default is used.
//   - x-migrations-table-quoted: boolean parsed with strconv.ParseBool. When
//     true, x-migrations-table must include surrounding double quotes.
//   - x-statement-timeout: integer statement timeout in milliseconds. Malformed
//     integers reject the URL.
//   - x-multi-statement: boolean parsed with strconv.ParseBool.
//   - x-multi-statement-max-size: integer byte limit for multi-statement
//     migrations. Empty or non-positive values use the upstream driver default.
//
// Malformed booleans or integers return parsing errors. An invalid quoted table
// configuration returns [ErrInvalidMigrationsTable].
func ParseConfig(u *url.URL) (*Config, error) {
	query := u.Query()
	migrationsTable, migrationsTableQuoted, err := parseMigrationsTableOptions(query)
	if err != nil {
		return nil, err
	}

	statementTimeout, err := parseIntOption(query, "x-statement-timeout")
	if err != nil {
		return nil, err
	}

	multiStatementMaxSize, err := parseMultiStatementMaxSize(query)
	if err != nil {
		return nil, err
	}

	multiStatementEnabled, err := parseBoolOption(query, "x-multi-statement")
	if err != nil {
		return nil, fmt.Errorf("unable to parse option x-multi-statement: %w", err)
	}

	config := &Config{
		MigrationsTable:       migrationsTable,
		MigrationsTableQuoted: migrationsTableQuoted,
		StatementTimeout:      (time.Duration(statementTimeout) * time.Millisecond).Duration(),
		MultiStatementEnabled: multiStatementEnabled,
		MultiStatementMaxSize: multiStatementMaxSize,
	}

	return config, nil
}

// WithInstance creates a migrate database driver from db and cfg.
func WithInstance(db *sql.DB, cfg *Config) (database.Driver, error) {
	return pgx.WithInstance(db, cfg)
}

func parseMigrationsTableOptions(query url.Values) (string, bool, error) {
	migrationsTable := query.Get("x-migrations-table")

	migrationsTableQuoted, err := parseBoolOption(query, "x-migrations-table-quoted")
	if err != nil {
		return "", false, fmt.Errorf("unable to parse option x-migrations-table-quoted: %w", err)
	}
	if migrationsTable != "" && migrationsTableQuoted && (migrationsTable[0] != '"' || migrationsTable[len(migrationsTable)-1] != '"') {
		return "", false, fmt.Errorf(
			"%w: x-migrations-table must be quoted when x-migrations-table-quoted is enabled, current value is %q",
			ErrInvalidMigrationsTable, migrationsTable,
		)
	}

	return migrationsTable, migrationsTableQuoted, nil
}

func parseMultiStatementMaxSize(query url.Values) (int, error) {
	value := query.Get("x-multi-statement-max-size")
	if value == "" {
		return pgx.DefaultMultiStatementMaxSize, nil
	}

	size, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	if size <= 0 {
		return pgx.DefaultMultiStatementMaxSize, nil
	}

	return size, nil
}

func parseBoolOption(query url.Values, name string) (bool, error) {
	value := query.Get(name)
	if value == "" {
		return false, nil
	}

	return strconv.ParseBool(value)
}

func parseIntOption(query url.Values, name string) (int, error) {
	value := query.Get(name)
	if value == "" {
		return 0, nil
	}

	return strconv.Atoi(value)
}
