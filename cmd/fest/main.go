package main

import (
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"

	h "github.com/upframe/fest/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/upframe/fest"
	"github.com/upframe/fest/email"
	"github.com/upframe/fest/mysql"
)

type config struct {
	Development bool   `json:"Development"`
	Key         string `json:"Key"`
	InviteOnly  bool   `json:"InviteOnly"`
	BaseInvites int    `json:"BaseInvites"`
	Database    struct {
		User     string `json:"User"`
		Password string `json:"Password"`
		Host     string `json:"Host"`
		Port     string `json:"Port"`
		Name     string `json:"Name"`
	} `json:"Database"`
	SMTP struct {
		User     string `json:"User"`
		Password string `json:"Password"`
		Host     string `json:"Host"`
		Port     string `json:"Port"`
	} `json:"SMTP"`
	PayPal struct {
		Client string `json:"Client"`
		Secret string `json:"Secret"`
	} `json:"PayPal"`
}

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

	c := &config{}
	// then config file settings

	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&c); err != nil {
		log.Fatal("parsing config file", err.Error())
	}

	fest.BaseInvites = c.BaseInvites
	fest.InviteOnly = c.InviteOnly

	fest.BaseAddress = "http://localhost"
	fest.Templates = "_assets/templates/"
	email.Templates = "_assets/templates/email/"

	// Creates the new cookie session;
	fest.Store = sessions.NewCookieStore([]byte(c.Key))
	fest.Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 3,
		Secure:   !c.Development,
		HttpOnly: true,
	}

	fest.Store.Options.Domain = "localhost"

	// Configures the email
	email.InitSMTP(c.SMTP.User, c.SMTP.Password, c.SMTP.Host, c.SMTP.Port)

	// Connects to the database and checks for an error
	db, err := mysql.InitDB(
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)

	if err != nil {
		panic(err)
	}

	// Configures PayPal
	if fest.InitPayPal(c.PayPal.Client, c.PayPal.Secret, c.Development) != nil {
		panic(err)
	}

	s := &fest.Services{
		User:      &mysql.UserService{DB: db},
		Link:      &mysql.LinkService{DB: db},
		Product:   &mysql.ProductService{DB: db},
		Promocode: &mysql.PromocodeService{DB: db},
		// Order: &mysql.OrderService{}
	}

	r := mux.NewRouter()

	r.Handle("/", &h.IndexHandler{
		Services: s,
	})

	r.Handle("/login", &h.LoginHandler{
		Services: s,
	})

	r.Handle("/register", &h.RegisterHandler{
		Services: s,
	})

	r.Handle("/reset", &h.ResetHandler{
		Services: s,
	})

	r.Handle("/settings", h.MustLogin(&h.SettingsHandler{
		Services: s,
	}))

	r.Handle("/store", &h.StoreHandler{
		Services: s,
	})

	r.Handle("/cart", h.MustLogin(&h.CartHandler{
		Services: s,
	}))

	r.Handle("/cart/{id}", h.MustLogin(&h.CartItemHandler{
		Services: s,
	}))

	r.Handle("/settings/deactivate", &h.DeactivateHandler{
		Services: s,
	})

	r.Handle("/checkout/cancel", &h.CheckoutCancelHandler{Services: s})
	/* r.Handle("/checkout/confirm", h.MustLogin(&h.CheckoutConfirmHandler{
		Services: s,
	})) */

	r.Handle("/checkout", h.MustLogin(&h.CheckoutHandler{
		Services: s,
	}))

	r.Handle("/coupon/validate", h.MustLogin(&h.ValidatePromocodeHandler{
		Services: s,
	}))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("_assets/static/"))))

	api := r.PathPrefix("/api").Subrouter()
	api.NotFoundHandler = &h.NotFoundAPI{}
	api.Handle("/promocode/{id:[0-9]+}", &h.PromocodeAPI{Services: s})

	/* api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/user/:id", GetUser).Methods("GET") */

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":80", h.InjectSession(r, s.User)))

	// TODO: check csrf things
}
