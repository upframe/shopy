package fest

import (
	"crypto/rand"
	"io"
	"log"
	"path/filepath"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/upframe/fest/email"
	"github.com/upframe/fest/models"
	"github.com/upframe/fest/pages"
)

var (
	// Store stores the session cookies and help us to handle them
	store *sessions.CookieStore
	// BaseAddress is the base URL to build URLs
	BaseAddress string
	// RootPath is the 'root' directive defined in Caddyfile
	RootPath string
	// TemplatesPath is where the templates are stored
	TemplatesPath string
)

func init() {
	// Regists the caddy middleware
	caddy.RegisterPlugin("fest", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	// Gets the base address
	cfg := httpserver.GetConfig(c)

	pages.BaseAddress = cfg.Addr.String()
	pages.Templates = filepath.Clean(cfg.Root+"/templates/") + string(filepath.Separator)
	email.Templates = pages.Templates + "email" + string(filepath.Separator)

	var (
		err                                    error
		smtpUser, smtpPass, smtpHost, smtpPort string
		dbUser, dbPass, dbHost, dbName         string
		paypalClient, paypalSecret             string
		dbPort                                 = "3306"
		development                            = false
		keyPairs                               [][]byte
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
			case "base_invites":
				if !c.NextArg() {
					return c.ArgErr()
				}

				pages.BaseInvites, err = strconv.Atoi(c.Val())
				if err != nil {
					return err
				}
			case "paypal_client":
				if !c.NextArg() {
					return c.ArgErr()
				}

				paypalClient = c.Val()
			case "paypal_secret":
				if !c.NextArg() {
					return c.ArgErr()
				}

				paypalSecret = c.Val()
			case "development":
				development = true
			case "invite_only":
				if !c.NextArg() {
					pages.InviteOnly = true
				} else {
					pages.InviteOnly, err = strconv.ParseBool(c.Val())
					if err != nil {
						return err
					}
				}
			}
		}
	}

	// Sets up the cookies
	if !development {
		// Generates 5 random key pairs to secure the cookies
		// NOTE: generating this at startup will automatically log out the
		// users when the server is rebooted
		keyPairs = [][]byte{}
		for i := 0; i < 5; i++ {
			keyPairs = append(keyPairs, make([]byte, 32))
			_, err = io.ReadFull(rand.Reader, keyPairs[i])
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		keyPairs = [][]byte{[]byte("HEY")}
	}

	// Creates the new cookie session;
	store = sessions.NewCookieStore(keyPairs...)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 3,
		Secure:   cfg.Addr.Scheme == "https",
		HttpOnly: true,
	}
	store.Options.Domain = cfg.Host()

	// Configures the email
	email.InitSMTP(smtpUser, smtpPass, smtpHost, smtpPort)

	// Connects to the database and checks for an error
	err = models.InitDB(dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		return err
	}

	// Configures PayPal
	if pages.InitPayPal(paypalClient, paypalSecret, development) != nil {
		return err
	}

	// Adds the middleware to Caddy
	mid := func(next httpserver.Handler) httpserver.Handler {
		return Upframe{
			Next: next,
			Root: cfg.Root,
		}
	}

	httpserver.GetConfig(c).AddMiddleware(mid)
	return nil
}
