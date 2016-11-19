package main

import (
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/upframe/fest"
	"github.com/upframe/fest/email"
	h "github.com/upframe/fest/http"
	"github.com/upframe/fest/mysql"
)

func main() {
	// Execute with all of the CPUs available
	runtime.GOMAXPROCS(runtime.NumCPU())

	path := "config.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	conf, err := fest.ConfigFile(path)
	if err != nil {
		panic(err)
	}

	// Connects to the database and checks for an error
	db, err := mysql.InitDB(
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name,
	)

	if err != nil {
		panic(err)
	}

	// Configures the email
	// TODO: Email... global or a struct inside config?
	email.Templates = conf.Assets + "templates/email/"
	email.InitSMTP(conf.SMTP.User, conf.SMTP.Password, conf.SMTP.Host, conf.SMTP.Port)

	// Configures PayPal
	paypal, err := fest.InitPayPal(conf.PayPal.Client, conf.PayPal.Secret, conf.Development)

	if err != nil {
		panic(err)
	}

	c := &fest.Config{
		Domain:         conf.Domain,
		Scheme:         conf.Scheme,
		Port:           strconv.Itoa(conf.Port),
		Assets:         conf.Assets,
		InviteOnly:     conf.InviteOnly,
		DefaultInvites: conf.DefaultInvites,
		BaseAddress:    conf.Scheme + "//" + conf.Domain,
		Templates:      conf.Assets + "templates/",
		Logger:         log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
		PayPal:         paypal,
		Store:          sessions.NewCookieStore([]byte(conf.Key)),
		Services: &fest.Services{
			User:      &mysql.UserService{DB: db},
			Link:      &mysql.LinkService{DB: db},
			Product:   &mysql.ProductService{DB: db},
			Promocode: &mysql.PromocodeService{DB: db},
			Order:     &mysql.OrderService{DB: db},
		},
	}

	// Define the Store options
	c.CookieStore = securecookie.New(conf.Key1, conf.Key2)

	h.Serve(c)
}
