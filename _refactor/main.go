package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/upframe/fest/pages"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", pages.IndexGET).Methods("GET")

	/* r.HandleFunc("/register", pages.RegisterGET).Methods("GET")
	r.HandleFunc("/register", pages.RegisterPOST).Methods("POST")

	r.HandleFunc("/login", pages.LoginGET).Methods("GET")
	r.HandleFunc("/login", pages.LoginPOST).Methods("POST")

	r.HandleFunc("/settings", pages.SettingsGET).Methods("GET")
	r.HandleFunc("/settings", pages.SettingsPUT).Methods("PUT") */

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("_assets/static/"))))

	/* api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/user/:id", GetUser).Methods("GET") */

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))

	// TODO: check csrf things
}

/*

	case r.URL.Path == "/settings/deactivate" && r.Method == http.MethodGet:
		return pages.DeactivateGET(w, r, s)
	case r.URL.Path == "/settings/deactivate" && r.Method == http.MethodPost:
		return pages.DeactivatePOST(w, r, s)
	case r.URL.Path == "/store" && r.Method == http.MethodGet:
		return pages.StoreGET(w, r, s)
	case r.URL.Path == "/cart" && r.Method == http.MethodGet:
		return pages.CartGET(w, r, s)
	case strings.HasPrefix(r.URL.Path, "/cart") && r.Method == http.MethodPost:
		return pages.CartPOST(w, r, s)
	case strings.HasPrefix(r.URL.Path, "/cart") && r.Method == http.MethodDelete:
		return pages.CartDELETE(w, r, s)
	case strings.HasPrefix(r.URL.Path, "/checkout"):
		return pages.Checkout(w, r, s)
	case r.URL.Path == "/coupon/validate" && r.Method == http.MethodPost:
		return pages.ValidatePromocode(w, r, s)
	case r.URL.Path == "/orders" && r.Method == http.MethodGet:
		return pages.OrdersGET(w, r, s)
	case r.URL.Path == "/reset" && r.Method == http.MethodGet:
		return pages.ResetGET(w, r, s)
	case r.URL.Path == "/reset" && r.Method == http.MethodPost:
		return pages.ResetPOST(w, r, s)
	case r.URL.Path == "/logout":
		return logout(w, r, s)
	}

	// Admin router: if the user is an admin and the page starts with /admin
	if s.IsAdmin() && strings.HasPrefix(r.URL.Path, "/admin") {
		if r.URL.Path == "/admin" && r.Method == http.MethodGet {
			return pages.RenderHTML(w, s, nil, "admin/home")
		}

		if strings.HasPrefix(r.URL.Path, "/admin/promocodes") {
			switch r.Method {
			case http.MethodGet:
				return pages.AdminPromocodesGET(w, r, s)
			case http.MethodPost:
				return pages.AdminPromocodesPOST(w, r)
			case http.MethodDelete:
				return pages.AdminPromocodesDELETE(w, r)
			case http.MethodPut:
				return pages.AdminPromocodesPUT(w, r)
			}
		}

		if strings.HasPrefix(r.URL.Path, "/admin/orders") {
			switch r.Method {
			case http.MethodGet:
				return pages.AdminOrdersGET(w, r, s)
			case http.MethodPost:
				return pages.AdminOrdersPOST(w, r)
			case http.MethodDelete:
				return pages.AdminOrdersDELETE(w, r)
			case http.MethodPut:
				return pages.AdminOrdersPUT(w, r)
			}
		}

		if strings.HasPrefix(r.URL.Path, "/admin/users") {
			switch r.Method {
			case http.MethodGet:
				return pages.AdminUsersGET(w, r, s)
			case http.MethodPost:
				return pages.AdminUsersPOST(w, r)
			case http.MethodDelete:
				return pages.AdminUsersDELETE(w, r)
			case http.MethodPut:
				return pages.AdminUsersPUT(w, r)
			}
		}

		if strings.HasPrefix(r.URL.Path, "/admin/products") {
			switch r.Method {
			case http.MethodGet:
				return pages.AdminProductsGET(w, r, s)
			case http.MethodPost:
				return pages.AdminProductsPOST(w, r)
			case http.MethodDelete:
				return pages.AdminProductsDELETE(w, r)
			case http.MethodPut:
				return pages.AdminProductsPUT(w, r)
			}
		}
	}

	// If the request doesn't match any route and it isn't a GET request
	// return a Status Not Implemented
	if r.Method != http.MethodGet {
		return http.StatusNotImplemented, nil
	}

	// Checks if there is a static template for this page. If so, show it!
	if _, err := os.Stat(filepath.Clean("templates/static" + r.URL.Path + ".tmpl")); err == nil {
		return pages.RenderHTML(w, nil, r.URL.Path)
	}

	// Return 404 Not Found for the rest
	return http.StatusNotFound, nil
}
*/
// logout resets the session values and saves the cookie
/* func logout(w http.ResponseWriter, r *http.Request) {
	s, _ := GetSession(w, r)
	// Reset the session values
	s.Values = map[interface{}]interface{}{}
	s.Values["IsLoggedIn"] = false

	// Saves the session and checks for error
	err := s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return http.StatusOK, nil
} */
