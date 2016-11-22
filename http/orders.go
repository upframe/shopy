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
	s := r.Context().Value("session").(*fest.Session)

	data, err := c.Services.Order.GetsWhere(0, 0, "ID", "User.ID", strconv.Itoa(s.User.ID))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return Render(w, c, s, data, "orders")
}

// OrderCancel ...
func OrderCancel(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	cart, err := c.Services.Cart.Get(w, r)
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

	if order.Status != fest.OrderWaitingPayment {
		return http.StatusNotFound, nil
	}

	if order.User.ID == s.User.ID {
		order.Status = fest.OrderCanceled
	} else {
		return http.StatusForbidden, nil
	}

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

	err = c.Services.Cart.Save(w, cart)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return Render(w, c, s, nil, "order-canceled")
}
