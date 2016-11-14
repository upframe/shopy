package http

import (
	"database/sql"
	"net/http"

	"github.com/upframe/fest"
)

// LoginHandler ...
type LoginHandler struct {
	UserService fest.UserService
	LinkService fest.LinkService
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		err  error
	)
	defer checkErrors(w, r, code, err)

	switch r.Method {
	case http.MethodGet:
		code, err = h.GET(w, r)
	case http.MethodPost:
		code, err = h.POST(w, r)
	default:
		code, err = http.StatusNotImplemented, nil
	}
}

// GET ...
func (h *LoginHandler) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	if s.IsLoggedIn() {
		return Redirect(w, r, "/")
	}

	return RenderHTML(w, s, nil, "login")
}

// POST ...
func (h *LoginHandler) POST(w http.ResponseWriter, r *http.Request) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	if s.IsLoggedIn() {
		return http.StatusBadRequest, fest.ErrAlreadyLoggedIn
	}

	if r.Header.Get("Resend") == "true" {
		// Obtains the user and checks for errors
		user, err := h.UserService.GetByEmail(r.Header.Get("email"))
		if err == sql.ErrNoRows {
			return http.StatusNotFound, err
		}

		if err != nil {
			return http.StatusInternalServerError, err
		}

		return confirmationEmail(h.LinkService, user)
	}

	// Parses the form and checks for errors
	err := r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Obtains the email and the password
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Checks if they're blank or not
	if email == "" || password == "" {
		return http.StatusBadRequest, nil
	}

	// Obtains the user and checks for errors
	user, err := h.UserService.GetByEmail(email)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Checks the password and checks for errors
	ok, err := user.CheckPassword(password)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !ok {
		return http.StatusUnauthorized, nil
	}

	// Checks if the user is confirmed
	if !user.Confirmed {
		return http.StatusFailedDependency, nil
	}

	// Checks if the user is deactivated
	if user.Deactivated {
		return http.StatusLocked, nil
	}

	// Sets the session cookie values
	s.Values["IsLoggedIn"] = true
	s.Values["UserID"] = user.ID
	s.Values["Cart"] = &fest.CartCookie{Products: map[int]int{}, Locked: false}
	s.Values["Order"] = &fest.OrderCookie{}

	// Saves the cookie and checks for errors
	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
