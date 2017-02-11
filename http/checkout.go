package http

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/logpacker/PayPal-Go-SDK"
	"github.com/bruhs/shopy"
)

// CheckoutConfirmGet ...
func CheckoutConfirmGet(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
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

	order, err := c.Services.Order.GetByPayPal(executeResult.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	order.PaymentStatus = paypalState(executeResult.State)
	err = c.Services.Order.Update(order, "PaymentStatus")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Saves the cookie and checks for errors
	err = c.Services.Cart.Reset(w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// TODO: send email with the invoice
	return Redirect(w, r, "/orders")
}

func paypalState(state string) int16 {
	switch state {
	case "created":
		return shopy.OrderPaymentWaiting
	case "approved":
		return shopy.OrderPaymentDone
	case "failed":
		return shopy.OrderPaymentFailed
	}

	return 0
}

// CheckoutGet ...
func CheckoutGet(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	s := r.Context().Value("session").(*shopy.Session)

	cart, err := c.Services.Cart.Get(w, r)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if len(cart.RawList) == 0 {
		return Redirect(w, r, "/cart")
	}

	err = cart.FillProducts(c.Services.Product)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	cart.Locked = true

	err = c.Services.Cart.Save(w, cart)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return Render(w, c, s, cart, "checkout")
}

// CheckoutPost ...
func CheckoutPost(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	s := r.Context().Value("session").(*shopy.Session)

	// Parses the form and checks for errors
	err := r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	cart, err := c.Services.Cart.Get(w, r)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = cart.FillProducts(c.Services.Product)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	order := &shopy.Order{
		User:          &shopy.User{ID: s.User.ID},
		PaymentStatus: shopy.OrderPaymentWaiting,
		Products:      []*shopy.OrderProduct{},
		Value:         cart.GetTotal(),
	}

	for _, product := range cart.Products {
		order.Products = append(order.Products, &shopy.OrderProduct{
			ID:       product.ID,
			Quantity: product.Quantity,
		})
	}

	// Obtain the credits and discount them
	text := r.FormValue("credits")
	if len(text) == 0 {
		text = "0"
	}

	credits, err := strconv.Atoi(text)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if s.User.Credit < credits || credits > order.Value {
		return http.StatusBadRequest, nil
	}

	order.Value -= credits

	// Gets the promocode from the form
	promocode := r.FormValue("promocode")

	if promocode != "" {
		// Gets the promocode and checks for errors
		order.Promocode, err = c.Services.Promocode.GetByCode(promocode)
		if err == sql.ErrNoRows {
			return http.StatusNotFound, nil
		}

		if err != nil {
			return http.StatusInternalServerError, err
		}

		if time.Now().Unix() > order.Promocode.Expires.Unix() || order.Promocode.Used == order.Promocode.MaxUsage {
			return http.StatusGone, nil
		}

		if order.Promocode.Used < order.Promocode.MaxUsage {
			order.Promocode.Used++
			err = c.Services.Promocode.Update(order.Promocode, "Used")
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}

		var discount int

		if order.Promocode.Percentage {
			discount = (order.Promocode.Discount * order.Value) / 100
		} else {
			discount = order.Promocode.Discount
		}

		order.Value -= discount
	}

	_, err = c.PayPal.GetAccessToken()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	amount := paypalsdk.Amount{
		Total:    shopy.DisplayCents(order.Value),
		Currency: "EUR",
	}

	err = c.Services.Order.Create(order)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	p, err := c.PayPal.CreateDirectPaypalPayment(
		amount,
		c.BaseAddress+"/checkout/confirm",
		c.BaseAddress+"/orders/"+strconv.Itoa(order.ID)+"/cancel",
		"Shop at Upframe Fest",
	)

	if err != nil || p.ID == "" {
		return http.StatusInternalServerError, err
	}

	order.PayPal = p.ID
	err = c.Services.Order.Update(order, "PayPal")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return apiPrint(w, map[string]string{"Link": p.Links[1].Href})
}
