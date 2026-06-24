package http

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/api/v1/migrate"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

// Register exposes the v1 service methods through the HTTP RPC facade.
func Register(migrator *migrate.Migrator) {
	rpc.Route(v1.Service_Migrate_FullMethodName, migrateDatabase(migrator))
	rpc.Route(v1.Service_ApplyMigrations_FullMethodName, applyMigrations(migrator))
	rpc.Route(v1.Service_Status_FullMethodName, getStatus(migrator))
	rpc.Route(v1.Service_ListDatabases_FullMethodName, migrator.ListDatabases)
}

func setFailureHeaders(ctx context.Context, values diagnostics.Values) {
	header := meta.Response(ctx).Header()
	for key, value := range values {
		header.Set(key, value)
	}
}
