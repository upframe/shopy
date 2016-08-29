package upframe

import (
	"github.com/hacdias/upframe/models"
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

	// TODO: get variables from caddyfile
	err := models.InitDB("root", "root", "127.0.0.1", "3306", "upframe")
	if err != nil {
		return err
	}

	httpserver.GetConfig(c).AddMiddleware(mid)
	return nil
}
