package database

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/migrieren/internal/migrate/pgx"
	"github.com/alexfalkowski/migrieren/internal/migrate/url"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/metric"
)

// ErrInvalidURL is returned when the url is invalid.
var ErrInvalidURL = errors.New("database: invalid url")

// ErrUnsupportedDriver is returned when the driver is not supported.
var ErrUnsupportedDriver = errors.New("database: unsupported driver")

var telemetryAttrs = telemetry.WithAttributes(attributes.DBSystemNamePostgreSQL)

// Open opens a migrate database driver for databaseURL.
//
// When ctx has a deadline, pgx statement execution is capped to no more than
// the remaining request duration. A shorter configured statement timeout is
// preserved.
//
// Empty, malformed, and unsupported-scheme URLs are returned to the caller via
// the exported sentinel errors in this package. Driver-specific option parsing
// errors, such as invalid pgx query parameters, are returned as-is.
//
// Telemetry wiring is treated differently on purpose: failures from
// telemetry.Open or telemetry.RegisterDBStatsMetrics are considered
// process-level misconfiguration/invariant violations for this service, so this
// function fails fast via runtime.Must rather than degrading to a runtime
// migration error.
func Open(ctx context.Context, databaseURL string) (database.Driver, error) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return nil, ErrInvalidURL
	}

	switch u.Scheme {
	case "pgx5":
		return openPGX(ctx, u)
	default:
		return nil, ErrUnsupportedDriver
	}
}

func openPGX(ctx context.Context, u *url.URL) (database.Driver, error) {
	cfg, err := pgx.ParseConfig(u)
	if err != nil {
		return nil, err
	}

	if deadline, ok := ctx.Deadline(); ok {
		timeout := time.Until(deadline).Duration()
		if timeout <= 0 {
			timeout = time.Nanosecond.Duration()
		}
		if cfg.StatementTimeout == 0 || timeout < cfg.StatementTimeout {
			cfg.StatementTimeout = timeout
		}
	}

	db, err := telemetry.Open("pgx/v5", url.DatabaseURL(u), telemetryAttrs)
	if err != nil {
		return nil, err
	}

	reg, err := telemetry.RegisterDBStatsMetrics(db, telemetryAttrs)
	if err != nil {
		return nil, err
	}

	dbDriver, err := pgx.WithInstance(db, cfg)
	if err != nil {
		_ = reg.Unregister()
		_ = db.Close()

		return nil, err
	}

	return &instrumentedDriver{Driver: dbDriver, registration: reg}, nil
}

// Ping opens databaseURL and verifies that the database can be reached with ctx.
func Ping(ctx context.Context, databaseURL string) error {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return ErrInvalidURL
	}

	switch u.Scheme {
	case "pgx5":
		if _, err := pgx.ParseConfig(u); err != nil {
			return err
		}

		db, err := telemetry.Open("pgx/v5", url.DatabaseURL(u), telemetryAttrs)
		if err != nil {
			return err
		}
		defer db.Close()

		return db.PingContext(ctx)
	default:
		return ErrUnsupportedDriver
	}
}

type instrumentedDriver struct {
	database.Driver
	registration metric.Registration
}

func (d *instrumentedDriver) Close() error {
	return errors.Join(d.Driver.Close(), d.registration.Unregister())
}
