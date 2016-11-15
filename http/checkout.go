package http

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/logpacker/PayPal-Go-SDK"
	"github.com/upframe/fest"
)

// TODO:

// CheckoutCancelGet ...
func CheckoutCancelGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	cart := s.Values["Cart"].(fest.CartCookie)
	cart.Locked = false

	s.Values["Order"] = &fest.OrderCookie{}
	s.Values["Cart"] = cart

	err := s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return Redirect(w, r, "/cart")
}

// CheckoutConfirmGet ...
func CheckoutConfirmGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	paymentID := r.URL.Query().Get("paymentId")
	payerID := r.URL.Query().Get("PayerID")

	if paymentID == "" || payerID == "" {
		return http.StatusBadRequest, nil
	}

	_, err := c.PayPal.GetAccessToken()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	executeResult, err := c.PayPal.ExecuteApprovedPayment(paymentID, payerID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	cart, err := s.GetCart(c.Services.Product)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	order := s.Values["Order"].(fest.OrderCookie)

	o := &fest.Order{
		UserID:    s.User.ID,
		PayPal:    paymentID,
		Value:     order.Total,
		Status:    executeResult.State,
		Promocode: &fest.Promocode{},
		Products:  []*fest.OrderProduct{},
	}

	if order.Promocode.Code != "" {
		o.Promocode.ID = order.Promocode.ID
	}

	for _, product := range cart.Products {
		o.Products = append(o.Products, &fest.OrderProduct{
			ID:       product.ID,
			Quantity: product.Quantity,
		})
	}

	err = c.Services.Order.Create(o)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	s.Values["Cart"] = &fest.CartCookie{Products: map[int]int{}, Locked: false}
	s.Values["Order"] = &fest.OrderCookie{}

	// Saves the cookie and checks for errors
	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// TODO: send email with the invoice
	return Redirect(w, r, "/orders")
}

// CheckoutGet ...
func CheckoutGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	cart, err := s.GetCart(c.Services.Product)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	cartCookie := s.Values["Cart"].(fest.CartCookie)
	cartCookie.Locked = true
	s.Values["Cart"] = cartCookie

	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err

	}

	return RenderHTML(w, c, s, cart, "checkout")
}

// CheckoutPost ...
func CheckoutPost(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	// Parses the form and checks for errors
	err := r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	cart, err := s.GetCart(c.Services.Product)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	order := fest.OrderCookie{}
	order.Total = cart.GetTotal()

	// Obtain the credits and discount them
	credits := r.FormValue("credits")
	if len(credits) == 0 {
		credits = "0"
	}

	order.Credits, err = strconv.Atoi(credits)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if s.User.Credit < order.Credits || order.Credits > order.Total {
		return http.StatusBadRequest, nil
	}

	order.Total -= order.Credits

	// Gets the promocode
	promocode := r.FormValue("promocode")

	if promocode != "" {
		// Gets the promocode and checks for errors
		var promo *fest.Promocode
		promo, err = c.Services.Promocode.GetByCode(promocode)
		if err == sql.ErrNoRows {
			return http.StatusNotFound, nil
		}

		if err != nil {
			return http.StatusInternalServerError, err
		}

		order.Promocode.Code = promo.Code
		order.Promocode.ID = promo.ID
		if time.Now().Unix() > promo.Expires.Unix() {
			return http.StatusGone, nil
		}

		if promo.Percentage {
			order.Promocode.DiscountAmount = (promo.Discount * order.Total) / 100
		} else {
			order.Promocode.DiscountAmount = promo.Discount
		}

		order.Total -= order.Promocode.DiscountAmount
	}

	// Saves the cookie
	s.Values["Order"] = order
	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = c.PayPal.GetAccessToken()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	amount := paypalsdk.Amount{
		Total:    displayCents(order.Total),
		Currency: "EUR",
	}

	p, err := c.PayPal.CreateDirectPaypalPayment(
		amount,
		c.BaseAddress+"/checkout/confirm",
		c.BaseAddress+"/checkout/cancel",
		"Shop at Upframe Fest",
	)

	if err != nil || p.ID == "" {
		return http.StatusInternalServerError, err
	}

	return Redirect(w, r, p.Links[1].Href)
}
