package http

import (
	"database/sql"
	"errors"
	"net/http"
	"net/mail"
	"time"

	"github.com/upframe/fest"
)

// RegisterGet ...
func RegisterGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)
	if s.Logged {
		return Redirect(w, r, "/")
	}

	if r.URL.Query().Get("confirm") != "" {
		link, err := c.Services.Link.Get(r.URL.Query().Get("confirm"))

		if err != nil || link.Used || link.Expires.Unix() < time.Now().Unix() || link.Path != "/register" {
			return Render(w, c, s, nil, "invalid-link")
		}

		user, err := c.Services.User.Get(link.User)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		user.Confirmed = true
		err = c.Services.User.Update(user, "Confirmed")
		if err != nil {
			return http.StatusInternalServerError, err
		}

		link.Used = true
		err = c.Services.Link.Update(link, "Used")
		if err != nil {
			return http.StatusInternalServerError, err
		}

		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return http.StatusOK, nil
	}

	if c.InviteOnly {
		// Gets the referrer user
		referrer, err := c.Services.User.GetByReferral(r.URL.Query().Get("ref"))

		// If the user doesn't exist show a page telling that registration
		// is invitation only
		if err == sql.ErrNoRows {
			return Render(w, c, s, nil, "register/invite")
		}

		if err != nil {
			return http.StatusInternalServerError, err
		}

		// If the user exists, but doesn't have invites, show that information
		if referrer.Invites < 1 {
			return Render(w, c, s, referrer, "register/gone")
		}
	}

	// Otherwise, show the registration page
	return Render(w, c, s, nil, "register/form")
}

// RegisterPost ...
func RegisterPost(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	if s.Logged {
		return http.StatusBadRequest, nil
	}

	// Parses the form and checks for errors
	err := r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Builds the user variable
	user := &fest.User{
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Email:     r.FormValue("email"),
		Invites:   c.DefaultInvites,
		Credit:    0,
		Confirmed: false,
		Referrer:  fest.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
	}

	if c.InviteOnly {
		// Gets the referrer user using the ?referral= option in the URL. If it doesn't
		// find the user, return a 403 Forbidden status
		var referrer *fest.User
		referrer, err = c.Services.User.GetByReferral(r.URL.Query().Get("ref"))
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
			err = c.Services.User.Update(referrer, "Invites")
		}()
	}

	// Checks if any of the fields is empty, if so, return a 400 Bad Request error
	if user.FirstName == "" || user.LastName == "" || user.Email == "" || r.FormValue("password") == "" {
		err = errors.New("First Name, Last Name, Email or Password is missing.")
		return http.StatusBadRequest, err
	}

	// Checks if there is already an user with this email. If there is,
	// return a 407 Conflict error.
	if is, _ := isExistentUser(c.Services.User, user.Email); is {
		err = errors.New("Email already registred.")
		return http.StatusConflict, err
	}

	// Generates a unique referral hash for this user
	user.Referral = fest.UniqueHash(user.Email)

	// Sets the password hash and salt for the user and checks for errors
	err = user.SetPassword(r.FormValue("password"))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Inserts the user into the database
	err = c.Services.User.Create(user)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Send the confirmation email. We return a StatusCreated even if we
	// get an error while sending the email. Why? Because the status response
	// is directed to the user CREATION. The user may ask for the resending
	// of the email if he needs to.
	_, err = confirmationEmail(c, user)
	return http.StatusCreated, err
}

func confirmationEmail(c *fest.Config, user *fest.User) (int, error) {
	// Sets the current time and expiration time of the confirmation email
	now := time.Now()
	expires := time.Now().Add(time.Hour * 24 * 20)

	link := &fest.Link{
		Path:    "/register",
		Hash:    fest.UniqueHash(user.Email),
		User:    user.ID,
		Used:    false,
		Time:    &now,
		Expires: &expires,
	}

	err := c.Services.Link.Create(link)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	data := make(map[string]interface{})
	data["Name"] = user.FirstName + " " + user.LastName
	data["Hash"] = link.Hash
	data["Host"] = c.BaseAddress

	email := &fest.Email{
		From: &mail.Address{
			Name: "Upframe",
		},
		To: &mail.Address{
			Name:    "",
			Address: user.Email,
		},
		Subject: "You're almost there",
	}

	err = c.Services.Email.UseTemplate(email, data, "confirmation")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = c.Services.Email.Send(email)
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
