package main

import (
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/mux"
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
		Services: &fest.Services{
			User:      &mysql.UserService{DB: db},
			Link:      &mysql.LinkService{DB: db},
			Product:   &mysql.ProductService{DB: db},
			Promocode: &mysql.PromocodeService{DB: db},
			// Order: &mysql.OrderService{}
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

	r := mux.NewRouter()

	// TODO:
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("_assets/static/"))))

	r.HandleFunc("/", h.Inject(h.IndexGet, c)).Methods("GET")
	r.HandleFunc("/login", h.Inject(h.LoginGet, c)).Methods("GET")
	r.HandleFunc("/login", h.Inject(h.LoginPost, c)).Methods("POST")

	r.HandleFunc("/register", h.Inject(h.RegisterGet, c)).Methods("GET")
	r.HandleFunc("/register", h.Inject(h.RegisterPost, c)).Methods("POST")

	r.HandleFunc("/settings", h.Inject(h.MustLogin(h.SettingsGet), c)).Methods("GET")
	r.HandleFunc("/settings", h.Inject(h.MustLogin(h.SettingsPost), c)).Methods("POST")

	r.HandleFunc("/settings/deactivate", h.Inject(h.MustLogin(h.DeactivateGet), c)).Methods("GET")
	r.HandleFunc("/settings/deactivate", h.Inject(h.MustLogin(h.DeactivatePost), c)).Methods("POST")

	r.HandleFunc("/store", h.Inject(h.StoreGet, c)).Methods("GET")

	r.HandleFunc("/cart", h.Inject(h.MustLogin(h.CartGet), c)).Methods("GET")

	r.HandleFunc("/cart/{id:[0-9]+}", h.Inject(h.MustLogin(h.CartItemPost), c)).Methods("POST")
	r.HandleFunc("/cart/{id:[0-9]+}", h.Inject(h.MustLogin(h.CartItemDelete), c)).Methods("DELETE")

	/*

		r.Handle("/checkout/cancel", &h.CheckoutCancelHandler{Services: s})
		/* r.Handle("/checkout/confirm", h.MustLogin(&h.CheckoutConfirmHandler{
			Services: s,
		})) */

	/* 	r.Handle("/checkout", h.MustLogin(&h.CheckoutHandler{
	   		Services: s,
	   	}))

	   	r.Handle("/coupon/validate", h.MustLogin(&h.ValidatePromocodeHandler{
	   		Services: s,
	   	}))

	   	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("_assets/static/"))))

	   	api := r.PathPrefix("/api").Subrouter()
	   	api.NotFoundHandler = &h.NotFoundAPI{}

	   	api.HandleFunc("/promocode/{id:[0-9]+}", h.APIPromocodeGET(c, s)).Methods("GET")

	   	/* api := r.PathPrefix("/api").Subrouter()
	   	api.HandleFunc("/user/:id", GetUser).Methods("GET") */

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":80", r))

	// TODO: check csrf things
}
