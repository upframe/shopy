package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/upframe/fest"
)

// Serve ...
func Serve(c *fest.Config) {
	r := mux.NewRouter()

	// TODO:
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("_assets/static/"))))

	r.HandleFunc("/", Inject(IndexGet, c)).Methods("GET")
	r.HandleFunc("/login", Inject(LoginGet, c)).Methods("GET")
	r.HandleFunc("/login", Inject(LoginPost, c)).Methods("POST")

	r.HandleFunc("/register", Inject(RegisterGet, c)).Methods("GET")
	r.HandleFunc("/register", Inject(RegisterPost, c)).Methods("POST")

	r.HandleFunc("/settings", Inject(MustLogin(SettingsGet), c)).Methods("GET")
	r.HandleFunc("/settings", Inject(MustLogin(SettingsPost), c)).Methods("POST")

	r.HandleFunc("/settings/deactivate", Inject(MustLogin(DeactivateGet), c)).Methods("GET")
	r.HandleFunc("/settings/deactivate", Inject(MustLogin(DeactivatePost), c)).Methods("POST")

	r.HandleFunc("/store", Inject(StoreGet, c)).Methods("GET")

	r.HandleFunc("/cart", Inject(MustLogin(CartGet), c)).Methods("GET")

	r.HandleFunc("/cart/{id:[0-9]+}", Inject(MustLogin(CartItemPost), c)).Methods("POST")
	r.HandleFunc("/cart/{id:[0-9]+}", Inject(MustLogin(CartItemDelete), c)).Methods("DELETE")

	/*

	   r.Handle("/checkout/cancel", &CheckoutCancelHandler{Services: s})
	   /* r.Handle("/checkout/confirm", MustLogin(&CheckoutConfirmHandler{
	       Services: s,
	   })) */

	/* 	r.Handle("/checkout", MustLogin(&CheckoutHandler{
	    Services: s,
	}))

	r.Handle("/coupon/validate", MustLogin(&ValidatePromocodeHandler{
	    Services: s,
	}))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("_assets/static/"))))

	api := r.PathPrefix("/api").Subrouter()
	api.NotFoundHandler = &NotFoundAPI{}

	api.HandleFunc("/promocode/{id:[0-9]+}", APIPromocodeGET(c, s)).Methods("GET")

	*/

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/order/{id:[0-9]+}", Inject(APIOrderGet, c)).Methods("GET")
	api.HandleFunc("/product/{id:[0-9]+}", Inject(APIProductGet, c)).Methods("GET")
	api.HandleFunc("/promocode/{id}", Inject(APIPromocodeGet, c)).Methods("GET")
	api.HandleFunc("/user/{id:[0-9]+}", Inject(APIUserGet, c)).Methods("GET")

	// TODO: POST; PUT AND DELETE

	// Bind to a port and pass our router in
	// TODO :check CSRF
	log.Fatal(http.ListenAndServe(":80", r))

}
