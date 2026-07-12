// Package migrate implements the migrieren.v1 migration contract shared by
// transports.
//
// The package accepts generated v1 request messages and returns generated v1
// response messages. It owns versioned API response construction and public
// request validation, then delegates configured database lookup, source/URL
// resolution, and core migration work to the API-facing migrate adapter.
//
// Transport packages remain responsible for registration, safe status-code
// mapping, and exposing diagnostics as the correct transport metadata
// mechanism, such as gRPC trailers or HTTP response headers.
package migrate
