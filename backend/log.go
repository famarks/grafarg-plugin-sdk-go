package backend

import (
	"github.com/famarks/grafarg-plugin-sdk-go/backend/log"
)

// Logger is the default logger instance. This can be used directly to log from
// your plugin to grafarg-server with calls like backend.Logger.Debug(...).
var Logger log.Logger = log.DefaultLogger
