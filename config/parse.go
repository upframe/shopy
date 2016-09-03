package config

import (
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/upframe/fest/models"
)

func init() {
	// Generates 5 random key pairs to secure the cookies
	// NOTE: generating this at startup will automatically log out the
	// users when the server is rebooted
	/* keyPairs := [][]byte{}

	   for i := 0; i < 5; i++ {
	       keyPairs = append(keyPairs, make([]byte, 32))
	       _, err := io.ReadFull(rand.Reader, keyPairs[i])
	       if err != nil {
	           log.Fatal(err)
	       }
	   } */

	// Creates the new cookie session; TODO: change this in production
	Store = sessions.NewCookieStore([]byte("HEY"))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 3, // 3 hours
		Secure:   false,    // TODO: Change this to true on the final website
		HttpOnly: true,
	}
}

// Parse parses the configuration from the Caddyfile to our program
func Parse(c *caddy.Controller) error {
	// Gets the base address
	cfg := httpserver.GetConfig(c)
	Store.Options.Domain = cfg.Host()

	RootPath = cfg.Root
	BaseAddress = cfg.Addr.String()
	TemplatesPath = filepath.Clean(cfg.Root+"/templates/") + "/"

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

				SMTPUser = c.Val()
			case "smtp_pass":
				if !c.NextArg() {
					return c.ArgErr()
				}

				SMTPPass = c.Val()
			case "smtp_host":
				if !c.NextArg() {
					return c.ArgErr()
				}

				SMTPHost = c.Val()
			case "smtp_port":
				if !c.NextArg() {
					return c.ArgErr()
				}

				SMTPPort = c.Val()
			}
		}
	}

	// Configures the email
	initSMTP()

	// Connects to the database and checks for an error
	err := initDB()
	if err != nil {
		return err
	}

	models.InitDB(db)

	return nil
}
