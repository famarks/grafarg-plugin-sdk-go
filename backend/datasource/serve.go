package datasource

import (
	"github.com/famarks/grafarg-plugin-sdk-go/backend"
)

// ServeOpts options for serving a data source plugin.
type ServeOpts struct {
	// CheckHealthHandler handler for health checks.
	// Optional to implement.
	backend.CheckHealthHandler

	// CallResourceHandler handler for resource calls.
	// Optional to implement.
	backend.CallResourceHandler

	// QueryDataHandler handler for data queries.
	// Required to implement.
	backend.QueryDataHandler

	// GRPCSettings settings for gPRC.
	GRPCSettings backend.GRPCSettings
}

// Serve starts serving the data source over gPRC.
func Serve(opts ServeOpts) error {
	return backend.Serve(backend.ServeOpts{
		CheckHealthHandler:  opts.CheckHealthHandler,
		CallResourceHandler: opts.CallResourceHandler,
		QueryDataHandler:    opts.QueryDataHandler,
		GRPCSettings:        opts.GRPCSettings,
	})
}
