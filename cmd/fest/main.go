package main

import (
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"

	handlers "github.com/upframe/fest/http"

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
	err = mysql.InitDB(
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
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create services.
	userService := &mysql.UserService{}
	linkService := &mysql.LinkService{}
	productService := &mysql.ProductService{}
	promoService := &mysql.PromocodeService{}
	orderService := &mysql.OrderService{}

	r := mux.NewRouter()

	r.Handle("/", &handlers.IndexHandler{
		UserService: userService,
	})

	r.Handle("/login", &handlers.LoginHandler{
		UserService: userService,
		LinkService: linkService,
	})

	r.Handle("/register", &handlers.RegisterHandler{
		UserService: userService,
		LinkService: linkService,
	})

	r.Handle("/reset", &handlers.ResetHandler{
		UserService: userService,
		LinkService: linkService,
	})

	r.Handle("/settings", &handlers.SettingsHandler{
		UserService: userService,
	})

	r.Handle("/store", &handlers.StoreHandler{
		UserService:     userService,
		ProductsService: productService,
	})

	r.Handle("/cart", &handlers.CartHandler{
		UserService:    userService,
		ProductService: productService,
	})

	r.Handle("/cart/{id}", &handlers.CartItemHandler{
		UserService:    userService,
		ProductService: productService,
	})

	r.Handle("/settings/deactivate", &handlers.DeactivateHandler{
		UserService: userService,
		LinkService: linkService,
	})

	r.Handle("/checkout/cancel", &handlers.CheckoutCancelHandler{
		UserService: userService,
	})
	r.Handle("/checkout/confirm", &handlers.CheckoutConfirmHandler{
		UserService:    userService,
		ProductService: productService,
		OrderService:   orderService,
	})

	r.Handle("/checkout", &handlers.CheckoutHandler{
		UserService:      userService,
		ProductService:   productService,
		PromocodeService: promoService,
	})

	r.Handle("/coupon/validate", &handlers.ValidatePromocodeHandler{
		UserService:      userService,
		PromocodeService: promoService,
	})

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("_assets/static/"))))

	/* api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/user/:id", GetUser).Methods("GET") */

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":80", r))

	// TODO: check csrf things
}
