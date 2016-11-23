package main

import (
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/gorilla/securecookie"
	"github.com/upframe/fest"
	h "github.com/upframe/fest/http"
	"github.com/upframe/fest/http/cookie"
	"github.com/upframe/fest/mysql"
	"github.com/upframe/fest/smtp"
)

func main() {
	// Execute with all of the CPUs available
	runtime.GOMAXPROCS(runtime.NumCPU())

	f, err := configFile("config.json")
	if err != nil {
		panic(err)
	}

	// Connects to the database and checks for an error
	db, err := mysql.InitDB(
		f.Database.User,
		f.Database.Password,
		f.Database.Host,
		f.Database.Port,
		f.Database.Name,
	)

	if err != nil {
		panic(err)
	}

	// figures PayPal
	paypal, err := fest.InitPayPal(f.PayPal.Client, f.PayPal.Secret, f.Development)

	if err != nil {
		panic(err)
	}

	email := smtp.InitSMTP(f.SMTP.User, f.SMTP.Password, f.SMTP.Host, f.SMTP.Port)
	email.TemplatesPath = f.Assets + "templates/email/"
	email.FromDefaultEmail = "noreply@upframe.xyz"

	store := securecookie.New([]byte(f.Key1), []byte(f.Key2))

	userService := &mysql.UserService{DB: db}

	c := &fest.Config{
		Domain:         f.Domain,
		Scheme:         f.Scheme,
		Port:           strconv.Itoa(f.Port),
		Assets:         f.Assets,
		InviteOnly:     f.InviteOnly,
		DefaultInvites: f.DefaultInvites,
		BaseAddress:    f.Scheme + "://" + f.Domain,
		Templates:      f.Assets + "templates/",
		PayPal:         paypal,
		CookieStore:    store,
		Services: &fest.Services{
			User:      userService,
			Link:      &mysql.LinkService{DB: db},
			Product:   &mysql.ProductService{DB: db},
			Promocode: &mysql.PromocodeService{DB: db},
			Order:     &mysql.OrderService{DB: db},
			Email:     email,
			Session: &cookie.SessionService{
				Store:       store,
				Secure:      f.Scheme == "https",
				UserService: userService,
			},
			Cart: &cookie.CartService{
				Store:  store,
				Secure: f.Scheme == "https",
			},
		},
	}

	if f.Errors == "" {
		f.Errors = "stdout"
	}

	if f.Errors == "stdout" {
		c.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	} else {
		var file *os.File
		file, err = os.OpenFile(f.Errors, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}

		defer file.Close()
		c.Logger = log.New(file, "", log.Ldate|log.Ltime)
	}

	h.Serve(c)
}
