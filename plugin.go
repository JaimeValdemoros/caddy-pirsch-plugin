package caddy_pirsch_plugin

import (
	"context"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	pirsch "github.com/pirsch-analytics/pirsch-go-sdk/v2/pkg"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(PirschPlugin{})
}

type PirschPlugin struct {
	ClientId     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	BaseURL      string `json:"base_url,omitempty"`

	logger *zap.Logger
	client *pirsch.Client
}

func (m PirschPlugin) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.pirsch",
		New: func() caddy.Module { return new(PirschPlugin) },
	}
}

func (m *PirschPlugin) Provision(ctx caddy.Context) (err error) {
	m.client = pirsch.NewClient(m.ClientId, m.ClientSecret, &pirsch.ClientConfig{
		BaseURL: strings.TrimSpace(m.BaseURL),
	})
	m.logger = ctx.Logger(m)
	return err
}

func (m *PirschPlugin) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	r2 := r.Clone(context.TODO())
	go func(r *http.Request) {
		options := new(pirsch.PageViewOptions)
		options.IP = r.Header.Get("X-Forwarded-For")
		if err := m.client.PageView(r, options); err != nil {
			m.logger.Error("failed sending page view to pirsch: %v", zap.Error(err))
		}
	}(r2)
	return next.ServeHTTP(w, r)
}

var _ caddyhttp.MiddlewareHandler = (*PirschPlugin)(nil)
