package http

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/upframe/fest"
)

// OrdersGet ...
func OrdersGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.SessionCookie)

	if !s.Logged {
		return Redirect(w, r, "/login")
	}

	data, err := c.Services.Order.GetsWhere(0, 0, "ID", "User.ID", strconv.Itoa(s.UserID))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return Render(w, c, s, data, "orders")
}

// OrderCancel ...
func OrderCancel(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	cart, err := ReadCartCookie(w, r, c)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	order, err := c.Services.Order.Get(id)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}

	order.Status = fest.OrderCanceled

	err = c.Services.Order.Update(order, "Status")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if order.Promocode != nil {
		if order.Promocode.Usage != -1 {
			order.Promocode.Usage++

			err = c.Services.Promocode.Update(order.Promocode, "Usage")
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}

	cart.Locked = false

	err = SetCartCookie(w, c, cart)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return Redirect(w, r, "/cart")
}
