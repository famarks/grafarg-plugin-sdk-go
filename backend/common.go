package backend

import (
	"context"
	"encoding/json"
	"time"

	"github.com/famarks/grafarg-plugin-sdk-go/backend/httpclient"
	"github.com/famarks/grafarg-plugin-sdk-go/backend/proxy"
	"github.com/famarks/grafarg-plugin-sdk-go/backend/useragent"
	"github.com/famarks/grafarg-plugin-sdk-go/internal/tenant"
)

const dataCustomOptionsKey = "grafargData"
const secureDataCustomOptionsKey = "grafargSecureData"

// User represents a Grafarg user.
type User struct {
	Login string
	Name  string
	Email string
	Role  string
}

// AppInstanceSettings represents settings for an app instance.
//
// In Grafarg an app instance is an app plugin of certain
// type that have been configured and enabled in a Grafarg organization.
type AppInstanceSettings struct {
	// JSONData repeats the properties at this level of the object (excluding DataSourceConfig), and also includes any
	// custom properties associated with the plugin config instance.
	JSONData json.RawMessage

	// DecryptedSecureJSONData contains key,value pairs where the encrypted configuration plugin instance in Grafarg
	// server have been decrypted before passing them to the plugin.
	DecryptedSecureJSONData map[string]string

	// Updated is the last time this plugin instance's configuration was updated.
	Updated time.Time
}

// HTTPClientOptions creates httpclient.Options based on settings.
func (s *AppInstanceSettings) HTTPClientOptions(_ context.Context) (httpclient.Options, error) {
	httpSettings, err := parseHTTPSettings(s.JSONData, s.DecryptedSecureJSONData)
	if err != nil {
		return httpclient.Options{}, err
	}

	opts := httpSettings.HTTPClientOptions()
	setCustomOptionsFromHTTPSettings(&opts, httpSettings)

	return opts, nil
}

// DataSourceInstanceSettings represents settings for a data source instance.
//
// In Grafarg a data source instance is a data source plugin of certain
// type that have been configured and created in a Grafarg organization.
type DataSourceInstanceSettings struct {
	// ID is the Grafarg assigned numeric identifier of the the data source instance.
	ID int64

	// UID is the Grafarg assigned string identifier of the the data source instance.
	UID string

	// Type is the unique identifier of the plugin that the request is for.
	// This should be the same value as PluginContext.PluginId.
	Type string

	// Name is the configured name of the data source instance.
	Name string

	// URL is the configured URL of a data source instance (e.g. the URL of an API endpoint).
	URL string

	// User is a configured user for a data source instance. This is not a Grafarg user, rather an arbitrary string.
	User string

	// Database is the configured database for a data source instance.
	// Only used by Elasticsearch and Influxdb.
	// Please use JSONData to store information related to database.
	Database string

	// BasicAuthEnabled indicates if this data source instance should use basic authentication.
	BasicAuthEnabled bool

	// BasicAuthUser is the configured user for basic authentication. (e.g. when a data source uses basic
	// authentication to connect to whatever API it fetches data from).
	BasicAuthUser string

	// JSONData contains the raw DataSourceConfig as JSON as stored by Grafarg server. It repeats the properties in
	// this object and includes custom properties.
	JSONData json.RawMessage

	// DecryptedSecureJSONData contains key,value pairs where the encrypted configuration in Grafarg server have been
	// decrypted before passing them to the plugin.
	DecryptedSecureJSONData map[string]string

	// Updated is the last time the configuration for the data source instance was updated.
	Updated time.Time
}

// HTTPClientOptions creates httpclient.Options based on settings.
func (s *DataSourceInstanceSettings) HTTPClientOptions(ctx context.Context) (httpclient.Options, error) {
	httpSettings, err := parseHTTPSettings(s.JSONData, s.DecryptedSecureJSONData)
	if err != nil {
		return httpclient.Options{}, err
	}

	if s.BasicAuthEnabled {
		httpSettings.BasicAuthEnabled = s.BasicAuthEnabled
		httpSettings.BasicAuthUser = s.BasicAuthUser
		httpSettings.BasicAuthPassword = s.DecryptedSecureJSONData["basicAuthPassword"]
	} else if s.User != "" {
		httpSettings.BasicAuthEnabled = true
		httpSettings.BasicAuthUser = s.User
		httpSettings.BasicAuthPassword = s.DecryptedSecureJSONData["password"]
	}

	opts := httpSettings.HTTPClientOptions()
	opts.Labels["datasource_name"] = s.Name
	opts.Labels["datasource_uid"] = s.UID
	opts.Labels["datasource_type"] = s.Type

	setCustomOptionsFromHTTPSettings(&opts, httpSettings)

	cfg := GrafargConfigFromContext(ctx)
	opts.ProxyOptions, err = s.ProxyOptions(cfg.proxy().clientCfg)
	if err != nil {
		return opts, err
	}

	return opts, nil
}

// PluginContext holds contextual information about a plugin request, such as
// Grafarg organization, user and plugin instance settings.
type PluginContext struct {
	// OrgID is The Grafarg organization identifier the request originating from.
	OrgID int64

	// PluginID is the unique identifier of the plugin that the request is for.
	PluginID string

	// PluginVersion is the version of the plugin that the request is for.
	PluginVersion string

	// User is the Grafarg user making the request.
	//
	// Will not be provided if Grafarg backend initiated the request,
	// for example when request is coming from Grafarg Alerting.
	User *User

	// AppInstanceSettings is the configured app instance settings.
	//
	// In Grafarg an app instance is an app plugin of certain
	// type that have been configured and enabled in a Grafarg organization.
	//
	// Will only be set if request targeting an app instance.
	AppInstanceSettings *AppInstanceSettings

	// DataSourceConfig is the configured data source instance
	// settings.
	//
	// In Grafarg a data source instance is a data source plugin of certain
	// type that have been configured and created in a Grafarg organization.
	//
	// Will only be set if request targeting a data source instance.
	DataSourceInstanceSettings *DataSourceInstanceSettings

	// GrafargConfig is the configuration settings provided by Grafarg.
	GrafargConfig *GrafargCfg

	// UserAgent is the user agent of the Grafarg server that initiated the gRPC request.
	// Will only be set if request is made from Grafarg v10.2.0 or later.
	UserAgent *useragent.UserAgent
}

func setCustomOptionsFromHTTPSettings(opts *httpclient.Options, httpSettings *HTTPSettings) {
	opts.CustomOptions = map[string]interface{}{}

	if httpSettings.JSONData != nil {
		opts.CustomOptions[dataCustomOptionsKey] = httpSettings.JSONData
	}

	if httpSettings.SecureJSONData != nil {
		opts.CustomOptions[secureDataCustomOptionsKey] = httpSettings.SecureJSONData
	}
}

// JSONDataFromHTTPClientOptions extracts JSON data from CustomOptions of httpclient.Options.
func JSONDataFromHTTPClientOptions(opts httpclient.Options) (res map[string]interface{}) {
	if opts.CustomOptions == nil {
		return
	}

	val, exists := opts.CustomOptions[dataCustomOptionsKey]
	if !exists {
		return
	}

	jsonData, ok := val.(map[string]interface{})
	if !ok {
		return
	}

	return jsonData
}

// SecureJSONDataFromHTTPClientOptions extracts secure JSON data from CustomOptions of httpclient.Options.
func SecureJSONDataFromHTTPClientOptions(opts httpclient.Options) (res map[string]string) {
	if opts.CustomOptions == nil {
		return
	}

	val, exists := opts.CustomOptions[secureDataCustomOptionsKey]
	if !exists {
		return
	}

	secureJSONData, ok := val.(map[string]string)
	if !ok {
		return
	}

	return secureJSONData
}

func propagateTenantIDIfPresent(ctx context.Context) context.Context {
	if tid, exists := tenant.IDFromIncomingGRPCContext(ctx); exists {
		ctx = tenant.WithTenant(ctx, tid)
	}
	return ctx
}

func (s *DataSourceInstanceSettings) ProxyOptions(clientCfg *proxy.ClientCfg) (*proxy.Options, error) {
	opts := &proxy.Options{}

	var dat map[string]interface{}
	if s.JSONData != nil {
		if err := json.Unmarshal(s.JSONData, &dat); err != nil {
			return nil, err
		}
	}

	opts.Enabled = proxy.SecureSocksProxyEnabledOnDS(dat)
	if !opts.Enabled {
		return nil, nil
	}

	opts.Auth = &proxy.AuthOptions{}
	opts.Timeouts = &proxy.TimeoutOptions{}
	if v, exists := dat["secureSocksProxyUsername"]; exists {
		opts.Auth.Username = v.(string)
	} else {
		// default username is the datasource uid
		opts.Auth.Username = s.UID
	}

	if v, exists := s.DecryptedSecureJSONData["secureSocksProxyPassword"]; exists {
		opts.Auth.Password = v
	}

	if v, exists := dat["timeout"]; exists {
		if iv, ok := v.(float64); ok {
			opts.Timeouts.Timeout = time.Duration(iv) * time.Second
		}
	} else {
		opts.Timeouts.Timeout = proxy.DefaultTimeoutOptions.Timeout
	}

	if v, exists := dat["keepAlive"]; exists {
		if iv, ok := v.(float64); ok {
			opts.Timeouts.KeepAlive = time.Duration(iv) * time.Second
		}
	} else {
		opts.Timeouts.KeepAlive = proxy.DefaultTimeoutOptions.KeepAlive
	}

	opts.ClientCfg = clientCfg

	return opts, nil
}