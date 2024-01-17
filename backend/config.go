package backend

import (
	"context"
	"strconv"
	"strings"

	"github.com/famarks/grafarg-plugin-sdk-go/backend/proxy"
	"github.com/famarks/grafarg-plugin-sdk-go/experimental/featuretoggles"
)

type configKey struct{}

// GrafargConfigFromContext returns Grafarg config from context.
func GrafargConfigFromContext(ctx context.Context) *GrafargCfg {
	v := ctx.Value(configKey{})
	if v == nil {
		return NewGrafargCfg(nil)
	}

	cfg := v.(*GrafargCfg)
	if cfg == nil {
		return NewGrafargCfg(nil)
	}

	return cfg
}

// WithGrafargConfig injects supplied Grafarg config into context.
func WithGrafargConfig(ctx context.Context, cfg *GrafargCfg) context.Context {
	ctx = context.WithValue(ctx, configKey{}, cfg)
	return ctx
}

type GrafargCfg struct {
	config map[string]string
}

func NewGrafargCfg(cfg map[string]string) *GrafargCfg {
	return &GrafargCfg{config: cfg}
}

func (c *GrafargCfg) Get(key string) string {
	return c.config[key]
}

func (c *GrafargCfg) FeatureToggles() FeatureToggles {
	features, exists := c.config[featuretoggles.EnabledFeatures]
	if !exists || features == "" {
		return FeatureToggles{}
	}

	fs := strings.Split(features, ",")
	enabledFeatures := make(map[string]struct{}, len(fs))
	for _, f := range fs {
		enabledFeatures[f] = struct{}{}
	}

	return FeatureToggles{
		enabled: enabledFeatures,
	}
}

func (c *GrafargCfg) Equal(c2 *GrafargCfg) bool {
	if c == nil && c2 == nil {
		return true
	}
	if c == nil || c2 == nil {
		return false
	}

	if len(c.config) != len(c2.config) {
		return false
	}
	for k, v1 := range c.config {
		if v2, ok := c2.config[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}

type FeatureToggles struct {
	// enabled is a set-like map of feature flags that are enabled.
	enabled map[string]struct{}
}

// IsEnabled returns true if feature f is contained in ft.enabled.
func (ft FeatureToggles) IsEnabled(f string) bool {
	_, exists := ft.enabled[f]
	return exists
}

type Proxy struct {
	clientCfg *proxy.ClientCfg
}

func (c *GrafargCfg) proxy() Proxy {
	if v, exists := c.config[proxy.PluginSecureSocksProxyEnabled]; exists && v == strconv.FormatBool(true) {
		return Proxy{
			clientCfg: &proxy.ClientCfg{
				ClientCert:   c.Get(proxy.PluginSecureSocksProxyClientCert),
				ClientKey:    c.Get(proxy.PluginSecureSocksProxyClientKey),
				RootCA:       c.Get(proxy.PluginSecureSocksProxyRootCACert),
				ProxyAddress: c.Get(proxy.PluginSecureSocksProxyProxyAddress),
				ServerName:   c.Get(proxy.PluginSecureSocksProxyServerName),
			},
		}
	}
	return Proxy{}
}
