package main

import (
	"encoding/gob"
	"encoding/json"
	"log"
	"os"
	"runtime"

	"github.com/gorilla/sessions"
	"github.com/upframe/fest"
	"github.com/upframe/fest/email"
	h "github.com/upframe/fest/http"
	"github.com/upframe/fest/mysql"
)

func init() {
	// Regist types so they can be used on Cookies
	gob.Register(fest.CartCookie{})
	gob.Register(fest.OrderCookie{})
}

func main() {
	// TODO: admin
	// TODO: api
	// TODO: clean

	runtime.GOMAXPROCS(runtime.NumCPU())

	file := &config{}
	// then config file settings

	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&file); err != nil {
		log.Fatal("parsing config file", err.Error())
	}

	// Connects to the database and checks for an error
	db, err := mysql.InitDB(
		file.Database.User,
		file.Database.Password,
		file.Database.Host,
		file.Database.Port,
		file.Database.Name,
	)

	if err != nil {
		panic(err)
	}

	c := &fest.Config{
		InviteOnly:     file.InviteOnly,
		DefaultInvites: file.DefaultInvites,
		BaseAddress:    "http://localhost",
		Templates:      "_assets/templates/",
		Logger:         log.New(os.Stdout, "", log.LstdFlags),
		Services: &fest.Services{
			User:      &mysql.UserService{DB: db},
			Link:      &mysql.LinkService{DB: db},
			Product:   &mysql.ProductService{DB: db},
			Promocode: &mysql.PromocodeService{DB: db},
			Order:     &mysql.OrderService{DB: db},
		},
	}

	email.Templates = "_assets/templates/email/"

	// Creates the new cookie session;
	c.Store = sessions.NewCookieStore([]byte(file.Key))
	c.Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 3,
		Secure:   !file.Development,
		HttpOnly: true,
		Domain:   "localhost",
	}

	// Configures the email
	email.InitSMTP(file.SMTP.User, file.SMTP.Password, file.SMTP.Host, file.SMTP.Port)

	// Configures PayPal
	paypal, err := fest.InitPayPal(file.PayPal.Client, file.PayPal.Secret, file.Development)

	if err != nil {
		panic(err)
	}

	c.PayPal = paypal

	h.Serve(c)
}
