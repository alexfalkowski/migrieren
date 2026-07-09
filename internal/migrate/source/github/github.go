package github

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/google/go-github/v72/github"
)

var (
	// ErrNoPassword is returned by Open when a github:// URL carries user info
	// but no password (personal access token).
	ErrNoPassword = errors.New("github source: no password provided")

	// ErrInvalidRepository is returned by Open when a github:// URL does not
	// contain a repository segment.
	ErrInvalidRepository = errors.New("github source: invalid repository")

	// ErrNotDirectory is returned by Open when the configured path resolves to a
	// file rather than a directory of migrations.
	ErrNotDirectory = errors.New("github source: path is not a directory")

	// ErrDuplicateMigration is returned by Open when the directory listing holds
	// more than one migration for the same version and direction.
	ErrDuplicateMigration = errors.New("github source: duplicate migration")
)

//nolint:gochecknoinits // source drivers register via source.Register in init; blank-import wiring depends on it.
func init() {
	source.Register("github", &Driver{})
}

// Driver is a golang-migrate source driver that reads migration files from a
// GitHub repository using github.com/google/go-github.
//
// The zero value is only valid as the registration placeholder whose Open
// method builds usable instances; all other methods assume an instance
// returned by Open.
type Driver struct {
	client     *github.Client
	options    *github.RepositoryContentGetOptions
	migrations *source.Migrations
	owner      string
	repository string
	path       string
}

// Open opens a GitHub migration source from a github:// URL of the form
//
//	github://[user:personal-access-token@]owner/repository[/path][#ref]
//
// The optional user info supplies a personal access token as the password; the
// user name is ignored. Without a token the repository is read
// unauthenticated. The optional fragment selects a git ref (branch, tag, or
// SHA) and defaults to the repository's default branch.
//
// Open lists the migration directory eagerly and returns [ErrNoPassword] when
// user info is present without a password, [ErrInvalidRepository] when no
// repository segment is present, [ErrNotDirectory] when the path resolves to a
// file, or the underlying go-github error when the listing fails.
//
// Open does not accept a context and is not request-scoped; callers that need
// bounded validation should avoid opening GitHub sources (see
// [github.com/alexfalkowski/migrieren/internal/migrate/source.Check]).
func (d *Driver) Open(rawURL string) (source.Driver, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	client := github.NewClient(nil)
	if parsed.User != nil {
		token, ok := parsed.User.Password()
		if !ok {
			return nil, ErrNoPassword
		}

		client = client.WithAuthToken(token)
	}

	repository, repoPath, _ := strings.Cut(strings.Trim(parsed.Path, "/"), "/")
	if repository == "" {
		return nil, ErrInvalidRepository
	}

	opened := &Driver{
		client:     client,
		options:    &github.RepositoryContentGetOptions{Ref: parsed.Fragment},
		migrations: source.NewMigrations(),
		owner:      parsed.Host,
		repository: repository,
		path:       repoPath,
	}

	if err := opened.readDirectory(); err != nil {
		return nil, err
	}

	return opened, nil
}

// Close releases resources held by the driver. The GitHub source holds no open
// handles, so Close always returns nil.
func (d *Driver) Close() error {
	return nil
}

// First returns the lowest migration version, or an os.ErrNotExist-wrapped
// error when the source has no migrations.
func (d *Driver) First() (uint, error) {
	version, ok := d.migrations.First()
	if !ok {
		return 0, &os.PathError{Op: "first", Path: d.path, Err: os.ErrNotExist}
	}

	return version, nil
}

// Prev returns the migration version preceding version, or an
// os.ErrNotExist-wrapped error when none exists.
func (d *Driver) Prev(version uint) (uint, error) {
	prev, ok := d.migrations.Prev(version)
	if !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("prev for version %d", version), Path: d.path, Err: os.ErrNotExist}
	}

	return prev, nil
}

// Next returns the migration version following version, or an
// os.ErrNotExist-wrapped error when none exists.
func (d *Driver) Next(version uint) (uint, error) {
	next, ok := d.migrations.Next(version)
	if !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("next for version %d", version), Path: d.path, Err: os.ErrNotExist}
	}

	return next, nil
}

// ReadUp returns the up migration body and its identifier for version, or an
// os.ErrNotExist-wrapped error when no up migration exists for version. The
// caller owns the returned reader and must close it.
func (d *Driver) ReadUp(version uint) (io.ReadCloser, string, error) {
	migration, ok := d.migrations.Up(version)
	if !ok {
		return nil, "", &os.PathError{Op: fmt.Sprintf("read up version %d", version), Path: d.path, Err: os.ErrNotExist}
	}

	return d.read(migration)
}

// ReadDown returns the down migration body and its identifier for version, or
// an os.ErrNotExist-wrapped error when no down migration exists for version.
// The caller owns the returned reader and must close it.
func (d *Driver) ReadDown(version uint) (io.ReadCloser, string, error) {
	migration, ok := d.migrations.Down(version)
	if !ok {
		return nil, "", &os.PathError{Op: fmt.Sprintf("read down version %d", version), Path: d.path, Err: os.ErrNotExist}
	}

	return d.read(migration)
}

func (d *Driver) readDirectory() error {
	file, directory, _, err := d.client.Repositories.GetContents(context.Background(), d.owner, d.repository, d.path, d.options)
	if err != nil {
		return err
	}

	if file != nil {
		return ErrNotDirectory
	}

	for _, content := range directory {
		migration, err := source.DefaultParse(content.GetName())
		if err != nil {
			continue // ignore files that are not migrations
		}

		if !d.migrations.Append(migration) {
			return errors.Prefix(content.GetName(), ErrDuplicateMigration)
		}
	}

	return nil
}

func (d *Driver) read(migration *source.Migration) (io.ReadCloser, string, error) {
	reader, _, err := d.client.Repositories.DownloadContents(
		context.Background(), d.owner, d.repository, path.Join(d.path, migration.Raw), d.options,
	)
	if err != nil {
		return nil, "", err
	}

	return reader, migration.Identifier, nil
}
