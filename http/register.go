package http

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/upframe/fest"
	"github.com/upframe/fest/crypto"
	"github.com/upframe/fest/email"
	"github.com/upframe/fest/utils/random"
)

// TODO: move to domain
var (
	BaseInvites = 0
	InviteOnly  = false
)

// RegisterHandler ...
type RegisterHandler struct {
	SessionService fest.SessionService
	LinkService    fest.LinkService
	UserService    fest.UserService
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		err  error
	)

	switch r.Method {
	case http.MethodGet:
		code, err = h.GET(w, r)
	case http.MethodPost:
		code, err = h.POST(w, r)
	default:
		code, err = http.StatusNotImplemented, nil
	}

	checkErrors(w, code, err)
}

// GET ...
func (h *RegisterHandler) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	s, err := h.SessionService.Session(w, r)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if s.IsLoggedIn() {
		return Redirect(w, r, "/")
	}

	if r.URL.Query().Get("confirm") != "" {
		var link *fest.Link
		link, err = h.LinkService.GetByHash(r.URL.Query().Get("confirm"))

		if err != nil || link.Used || link.Expires.Unix() < time.Now().Unix() || link.Path != "/register" {
			return RenderHTML(w, s, nil, "invalid-link")
		}

		user, err := h.UserService.Get(link.User)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		user.Confirmed = true
		err = h.UserService.Update(user, "Confirmed")
		if err != nil {
			return http.StatusInternalServerError, err
		}

		link.Used = true
		err = h.LinkService.Update(link, "Used")
		if err != nil {
			return http.StatusInternalServerError, err
		}

		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return http.StatusOK, nil
	}

	if InviteOnly {
		// Gets the referrer user
		referrer, err := h.UserService.GetByReferral(r.URL.Query().Get("ref"))

		// If the user doesn't exist show a page telling that registration
		// is invitation only
		if err != nil {
			log.Println(err)
			return RenderHTML(w, s, nil, "register/invite")
		}

		// If the user exists, but doesn't have invites, show that information
		if referrer.Invites < 1 {
			return RenderHTML(w, s, referrer, "register/gone")
		}
	}

	// Otherwise, show the registration page
	return RenderHTML(w, s, nil, "register/form")
}

// POST ...
func (h *RegisterHandler) POST(w http.ResponseWriter, r *http.Request) (int, error) {
	s, err := h.SessionService.Session(w, r)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if s.IsLoggedIn() {
		return http.StatusBadRequest, nil
	}

	// Parses the form and checks for errors
	err = r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Builds the user variable
	user := &fest.User{
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Email:     r.FormValue("email"),
		Invites:   BaseInvites,
		Credit:    0,
		Confirmed: false,
		Referrer:  fest.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
	}

	if InviteOnly {
		// Gets the referrer user using the ?referral= option in the URL. If it doesn't
		// find the user, return a 403 Forbidden status
		var referrer *fest.User
		referrer, err = h.UserService.GetByReferral(r.URL.Query().Get("ref"))
		if err != nil {
			return http.StatusForbidden, err
		}

		// Checks if the referrer still has invites available! This is important! If
		// it doesn't, return a 410 Status Gone
		if referrer.Invites < 1 {
			return http.StatusGone, nil
		}

		user.Referrer = fest.NullInt64{
			NullInt64: sql.NullInt64{
				Int64: int64(referrer.ID),
				Valid: true,
			},
		}

		// This is the last thing to do. If there is a problem removing one invite,
		// who cares? We just need to make sure that everything is logged!
		defer func() {
			if err != nil {
				return
			}

			// Decrement one value from the referrer invites number and updates it in
			// the database and checks for errors. In this case, if there is an error
			// we will keep the registration going because it's no user-fault and the
			// refferer had invites available.
			referrer.Invites--
			err = h.UserService.Update(referrer, "Invites")
		}()
	}

	// Checks if any of the fields is empty, if so, return a 400 Bad Request error
	if user.FirstName == "" || user.LastName == "" || user.Email == "" || r.FormValue("password") == "" {
		err = errors.New("First Name, Last Name, Email or Password is missing.")
		return http.StatusBadRequest, err
	}

	// Checks if there is already an user with this email. If there is,
	// return a 407 Conflict error.
	if is, _ := isExistentUser(h.UserService, user.Email); is {
		err = errors.New("Email already registred.")
		return http.StatusConflict, err
	}

	// Generates a unique referral hash for this user
	user.Referral = random.UniqueHash(user.Email)

	// Sets the password hash and salt for the user and checks for errors
	err = crypto.SetPassword(user, r.FormValue("password"))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Inserts the user into the database
	user.ID, err = h.UserService.Create(user)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Send the confirmation email. We return a StatusCreated even if we
	// get an error while sending the email. Why? Because the status response
	// is directed to the user CREATION. The user may ask for the resending
	// of the email if he needs to.
	_, err = confirmationEmail(h.LinkService, user)
	return http.StatusCreated, err
}

func confirmationEmail(s fest.LinkService, user *fest.User) (int, error) {
	// Sets the current time and expiration time of the confirmation email
	now := time.Now()
	expires := time.Now().Add(time.Hour * 24 * 20)

	link := &fest.Link{
		Path:    "/register",
		Hash:    random.UniqueHash(user.Email),
		User:    user.ID,
		Used:    false,
		Time:    &now,
		Expires: &expires,
	}

	err := s.Create(link)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	data := make(map[string]interface{})
	data["Name"] = user.FirstName + " " + user.LastName
	data["Hash"] = link.Hash
	data["Host"] = BaseAddress

	email := &email.Email{
		From: &mail.Address{
			Name:    "Upframe",
			Address: email.FromDefaultEmail,
		},
		To: &mail.Address{
			Name:    "",
			Address: user.Email,
		},
		Subject: "You're almost there",
	}

	err = email.UseTemplate("confirmation", data)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = email.Send()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// isExistentUser checks if there is an user with the specified email
// and returns true and nil if the user exists and there is no error
func isExistentUser(s fest.UserService, mail string) (bool, error) {
	// Fetches the user from the database and checks for errors
	user, err := s.GetByEmail(mail)
	if err != nil {
		return false, err
	}

	// Checks if the user ID is different from 0, which means that it is valid
	// if so, returns true and nil
	return (user.ID != 0), nil
}
