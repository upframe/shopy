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

	r.HandleFunc("/orders", Inject(MustLogin(OrdersGet), c)).Methods("GET")

	r.HandleFunc("/checkout", Inject(MustLogin(CheckoutGet), c)).Methods("GET")
	r.HandleFunc("/checkout", Inject(MustLogin(CheckoutPost), c)).Methods("POST")
	r.HandleFunc("/checkout/cancel", Inject(MustLogin(CheckoutCancelGet), c)).Methods("GET")
	r.HandleFunc("/checkout/confirm", Inject(MustLogin(CheckoutConfirmGet), c)).Methods("GET")

	r.HandleFunc("/logout", Inject(logout, c))

	api := r.PathPrefix("/api").Subrouter()

	// Users can only access their own orders and their own user information. Admins
	// can access everything.
	api.HandleFunc("/order/{id:[0-9]+}", Inject(MustLogin(APIOrderGet), c)).Methods("GET")
	api.HandleFunc("/product/{id:[0-9]+}", Inject(MustAdmin(APIProductGet), c)).Methods("GET")
	api.HandleFunc("/promocode/{id}", Inject(MustLogin(APIPromocodeGet), c)).Methods("GET")
	api.HandleFunc("/user/{id:[0-9]+}", Inject(MustLogin(APIUserGet), c)).Methods("GET")

	api.HandleFunc("/order", Inject(MustAdmin(APIOrderPost), c)).Methods("POST")
	api.HandleFunc("/product", Inject(APIProductPost, c)).Methods("POST")
	api.HandleFunc("/promocode", Inject(MustAdmin(APIPromocodePost), c)).Methods("POST")
	api.HandleFunc("/user", Inject(MustAdmin(APIUserPost), c)).Methods("POST")

	// TODO: PUT

	api.HandleFunc("/order/{id:[0-9]+}", Inject(MustAdmin(APIOrderDelete), c)).Methods("DELETE")
	api.HandleFunc("/product/{id:[0-9]+}", Inject(APIProductDelete, c)).Methods("DELETE")
	api.HandleFunc("/promocode/{id:[0-9]+}", Inject(MustAdmin(APIPromocodeDelete), c)).Methods("DELETE")
	api.HandleFunc("/user/{id:[0-9]+}", Inject(MustAdmin(APIUserDelete), c)).Methods("DELETE")

	// Bind to a port and pass our router in
	// TODO :check CSRF
	log.Fatal(http.ListenAndServe(":80", r))
}

// logout resets the session values and saves the cookie
func logout(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	// Reset the session values
	s.Values = map[interface{}]interface{}{}
	s.Values["IsLoggedIn"] = false

	// Saves the session and checks for error
	err := s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return Redirect(w, r, "/")
}
