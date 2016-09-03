package upframe

import (
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/upframe/fest/config"
)

func init() {
	// Regists the caddy middleware
	caddy.RegisterPlugin("fest", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	err := config.Parse(c)
	if err != nil {
		return err
	}

	// Adds the middleware to Caddy
	mid := func(next httpserver.Handler) httpserver.Handler {
		return Upframe{
			Next: next,
		}
	}

	httpserver.GetConfig(c).AddMiddleware(mid)
	return nil
}
