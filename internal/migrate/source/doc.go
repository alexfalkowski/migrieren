// Package source opens migrate source drivers and registers supported schemes.
//
// Driver registration is performed via side-effect imports in this package.
// Supported source schemes are currently:
//   - file://
//   - github://
//
// [Open] delegates directly to github.com/golang-migrate/migrate/v4/source.Open.
package source
