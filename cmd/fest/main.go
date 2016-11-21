package main

import (
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/gorilla/securecookie"
	"github.com/upframe/fest"
	h "github.com/upframe/fest/http"
	"github.com/upframe/fest/mysql"
	"github.com/upframe/fest/smtp"
)

func main() {
	// Execute with all of the CPUs available
	runtime.GOMAXPROCS(runtime.NumCPU())

	path := "config.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	conf, err := configFile(path)
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

	// Configures PayPal
	paypal, err := fest.InitPayPal(conf.PayPal.Client, conf.PayPal.Secret, conf.Development)

	if err != nil {
		panic(err)
	}

	email := smtp.InitSMTP(conf.SMTP.User, conf.SMTP.Password, conf.SMTP.Host, conf.SMTP.Port)
	email.TemplatesPath = conf.Assets + "templates/email/"
	email.FromDefaultEmail = "noreply@upframe.xyz"

	c := &fest.Config{
		Domain:         conf.Domain,
		Scheme:         conf.Scheme,
		Port:           strconv.Itoa(conf.Port),
		Assets:         conf.Assets,
		InviteOnly:     conf.InviteOnly,
		DefaultInvites: conf.DefaultInvites,
		BaseAddress:    conf.Scheme + "://" + conf.Domain,
		Templates:      conf.Assets + "templates/",
		PayPal:         paypal,
		CookieStore:    securecookie.New([]byte(conf.Key1), []byte(conf.Key2)),
		Services: &fest.Services{
			User:      &mysql.UserService{DB: db},
			Link:      &mysql.LinkService{DB: db},
			Product:   &mysql.ProductService{DB: db},
			Promocode: &mysql.PromocodeService{DB: db},
			Order:     &mysql.OrderService{DB: db},
			Email:     email,
		},
	}

	if conf.Errors == "" {
		conf.Errors = "stdout"
	}

	if conf.Errors == "stdout" {
		c.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		var file *os.File
		file, err = os.OpenFile(conf.Errors, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}

		defer file.Close()
		c.Logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	h.Serve(c)
}
