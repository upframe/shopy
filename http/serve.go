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

	r.NotFoundHandler = &notFoundHandler{Config: c}

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(c.Assets+"static/"))))

	r.HandleFunc("/", Inject(StaticHandler("index"), c)).Methods("GET")
	r.HandleFunc("/login", Inject(LoginGet, c)).Methods("GET")
	r.HandleFunc("/login", Inject(LoginPost, c)).Methods("POST")

	r.HandleFunc("/register", Inject(RegisterGet, c)).Methods("GET")
	r.HandleFunc("/register", Inject(RegisterPost, c)).Methods("POST")

	r.HandleFunc("/reset", Inject(ResetGet, c)).Methods("GET")
	r.HandleFunc("/reset", Inject(ResetPost, c)).Methods("POST")

	r.HandleFunc("/settings", Inject(MustLogin(SettingsGet), c)).Methods("GET")

	r.HandleFunc("/settings/deactivate", Inject(MustLogin(DeactivateGet), c)).Methods("GET")
	r.HandleFunc("/settings/deactivate", Inject(MustLogin(DeactivatePost), c)).Methods("POST")

	r.HandleFunc("/store", Inject(StoreGet, c)).Methods("GET")

	// If you can Register without having an invite, you will be able to start adding
	// items to the cart without being logged in. Then, you log in or create an account
	// and finish the checkout.
	// TODO: this should be checked directly on the functions...
	if c.InviteOnly {
		r.HandleFunc("/cart", Inject(MustLogin(CartGet), c)).Methods("GET")
		r.HandleFunc("/cart/{id:[0-9]+}", Inject(MustLogin(CartItemPost), c)).Methods("POST")
		r.HandleFunc("/cart/{id:[0-9]+}", Inject(MustLogin(CartItemDelete), c)).Methods("DELETE")
	} else {
		r.HandleFunc("/cart", Inject(CartGet, c)).Methods("GET")
		r.HandleFunc("/cart/{id:[0-9]+}", Inject(CartItemPost, c)).Methods("POST")
		r.HandleFunc("/cart/{id:[0-9]+}", Inject(CartItemDelete, c)).Methods("DELETE")
	}

	r.HandleFunc("/orders", Inject(MustLogin(OrdersGet), c)).Methods("GET")
	r.HandleFunc("/orders/{id:[0-9]+}/cancel", Inject(MustLogin(OrderCancel), c)).Methods("GET")

	r.HandleFunc("/checkout", Inject(MustLogin(CheckoutGet), c)).Methods("GET")
	r.HandleFunc("/checkout", Inject(MustLogin(CheckoutPost), c)).Methods("POST")
	r.HandleFunc("/checkout/confirm", Inject(MustLogin(CheckoutConfirmGet), c)).Methods("GET")

	r.HandleFunc("/logout", Inject(logout, c))

	// Users can only access their own orders and their own user information. Admins
	// can access everything.
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/orders/{id:[0-9]+}", Inject(MustLogin(APIOrderGet), c)).Methods("GET")
	api.HandleFunc("/products/{id:[0-9]+}", Inject(MustAdmin(APIProductGet), c)).Methods("GET")
	api.HandleFunc("/promocodes/{id}", Inject(MustLogin(APIPromocodeGet), c)).Methods("GET")
	api.HandleFunc("/users/{id:[0-9]+}", Inject(MustLogin(APIUserGet), c)).Methods("GET")

	api.HandleFunc("/orders", Inject(MustAdmin(APIOrderPost), c)).Methods("POST")
	api.HandleFunc("/products", Inject(MustAdmin(APIProductPost), c)).Methods("POST")
	api.HandleFunc("/promocodes", Inject(MustAdmin(APIPromocodePost), c)).Methods("POST")
	api.HandleFunc("/users", Inject(MustAdmin(APIUserPost), c)).Methods("POST")

	api.HandleFunc("/orders/{id:[0-9]+}", Inject(MustAdmin(APIOrderPatch), c)).Methods("PATCH")
	api.HandleFunc("/products/{id:[0-9]+}", Inject(MustAdmin(APIProductPatch), c)).Methods("PATCH")
	api.HandleFunc("/promocodes/{id:[0-9]+}", Inject(MustAdmin(APIPromocodePatch), c)).Methods("PATCH")
	api.HandleFunc("/users/{id:[0-9]+}", Inject(MustLogin(APIUserPatch), c)).Methods("PATCH")

	api.HandleFunc("/orders/{id:[0-9]+}", Inject(MustAdmin(APIOrderDelete), c)).Methods("DELETE")
	api.HandleFunc("/products/{id:[0-9]+}", Inject(MustAdmin(APIProductDelete), c)).Methods("DELETE")
	api.HandleFunc("/promocodes/{id:[0-9]+}", Inject(MustAdmin(APIPromocodeDelete), c)).Methods("DELETE")
	api.HandleFunc("/users/{id:[0-9]+}", Inject(MustAdmin(APIUserDelete), c)).Methods("DELETE")
	api.HandleFunc("/users/current", Inject(MustLogin(APICurrentUser), c))

	r.HandleFunc("/admin", Inject(MustAdmin(AdminGet), c)).Methods("GET")

	admin := r.PathPrefix("/admin").Methods("GET").Subrouter()

	admin.HandleFunc("/{category:(?:products|promocodes|orders|users)}", Inject(MustAdmin(AdminRedirect), c))
	admin.HandleFunc("/{category:(?:products|promocodes|orders|users)}/new", Inject(MustAdmin(AdminNew), c))
	admin.HandleFunc("/{category:(?:products|promocodes|orders|users)}/last", Inject(MustAdmin(AdminListingLast), c))
	admin.HandleFunc("/{category:(?:products|promocodes|orders|users)}/{page:[0-9]+}", Inject(MustAdmin(AdminListing), c))

	// TODO :check CSRF
	log.Fatal(http.ListenAndServe(":"+c.Port, r))
}

type notFoundHandler struct {
	Config *fest.Config
}

func (h *notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Inject(func(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
		return http.StatusNotFound, nil
	}, h.Config)(w, r)
}

// logout resets the session values and saves the cookie
func logout(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {

	// Reset the session values
	err := SetSessionCookie(w, c, &fest.SessionCookie{})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return Redirect(w, r, "/")
}
