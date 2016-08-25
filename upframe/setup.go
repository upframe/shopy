package upframe

import (
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("upframe", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	mid := func(next httpserver.Handler) httpserver.Handler {
		return Upframe{
			Next: next,
		}
	}

	httpserver.GetConfig(c).AddMiddleware(mid)
	return nil
}
