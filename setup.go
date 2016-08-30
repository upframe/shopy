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
	// Initialize our pretty variables
	var (
		smtpUser, smtpPass, smtpHost, smtpPort string
		dbUser, dbPass, dbHost, dbPort, dbName string
	)

	// Gets the options from the Caddyfile
	for c.Next() {
		for c.NextBlock() {
			switch c.Val() {
			case "db_user":
				if !c.NextArg() {
					return c.ArgErr()
				}

				dbUser = c.Val()
			case "db_pass":
				if !c.NextArg() {
					return c.ArgErr()
				}

				dbPass = c.Val()
			case "db_host":
				if !c.NextArg() {
					return c.ArgErr()
				}

				dbHost = c.Val()
			case "db_port":
				if !c.NextArg() {
					return c.ArgErr()
				}

				dbPort = c.Val()
			case "db_name":
				if !c.NextArg() {
					return c.ArgErr()
				}

				dbName = c.Val()
			case "smtp_user":
				if !c.NextArg() {
					return c.ArgErr()
				}

				smtpUser = c.Val()
			case "smtp_pass":
				if !c.NextArg() {
					return c.ArgErr()
				}

				smtpPass = c.Val()
			case "smtp_host":
				if !c.NextArg() {
					return c.ArgErr()
				}

				smtpHost = c.Val()
			case "smtp_port":
				if !c.NextArg() {
					return c.ArgErr()
				}

				smtpPort = c.Val()
			}
		}
	}

	// Configures the email
	models.InitSMTP(smtpUser, smtpPass, smtpHost, smtpPort)

	// Connects to the database and checks for an error
	err := models.InitDB(dbUser, dbPass, dbHost, dbPort, dbName)
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
