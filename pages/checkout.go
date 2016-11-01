package pages

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/logpacker/PayPal-Go-SDK"
	"github.com/upframe/fest/models"
)

var c *paypalsdk.Client

// InitPayPal configures the paypal client variable
func InitPayPal(client, secret string, development bool) error {
	link := paypalsdk.APIBaseLive
	if development {
		link = paypalsdk.APIBaseSandBox
	}

	var err error
	c, err = paypalsdk.NewClient(client, secret, link)
	return err
}

// Checkout ...
func Checkout(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	if !s.IsLoggedIn() {
		if r.Method == http.MethodGet {
			return Redirect(w, r, "/login")
		}

		return http.StatusUnauthorized, errNotLoggedIn
	}

	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/checkout/confirm":
			return checkoutConfirmGET(w, r, s)
		case "/checkout/cancel":
			return checkoutCancelGET(w, r, s)
		case "/checkout":
			return checkoutGET(w, r, s)
		}

		return http.StatusNotFound, nil
	case http.MethodPost:
		switch r.URL.Path {
		case "/checkout":
			return checkoutPOST(w, r, s)
		}

		return http.StatusNotFound, nil
	}

	return http.StatusNotImplemented, nil
}

// CheckoutGET handles the GET request for /checkout page
func checkoutGET(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	cart, err := s.GetCart()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	cartCookie := s.Values["Cart"].(models.CartCookie)
	cartCookie.Locked = true
	s.Values["Cart"] = cartCookie

	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return RenderHTML(w, s, cart, "checkout")
}

func checkoutCancelGET(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	cart := s.Values["Cart"].(models.CartCookie)
	cart.Locked = false

	s.Values["Order"] = &models.OrderCookie{}
	s.Values["Cart"] = cart

	err := s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return Redirect(w, r, "/cart")
}

func checkoutConfirmGET(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	paymentID := r.URL.Query().Get("paymentId")
	payerID := r.URL.Query().Get("PayerID")

	if paymentID == "" || payerID == "" {
		return http.StatusBadRequest, nil
	}

	_, err := c.GetAccessToken()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	executeResult, err := c.ExecuteApprovedPayment(paymentID, payerID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	cart, err := s.GetCart()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	order := s.Values["Order"].(models.OrderCookie)

	o := &models.Order{
		UserID:      s.User.ID,
		PayPalID:    paymentID,
		Value:       order.Total,
		Status:      executeResult.State,
		PromocodeID: models.NullInt64JSON{},
	}
	if order.Promocode.Code == "" {
		o.PromocodeID.Valid = false
	} else {
		o.PromocodeID.Valid = true
		o.PromocodeID.Int64 = int64(order.Promocode.ID)
	}

	orderID, err := o.Insert()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	for _, product := range cart.Products {
		op := models.OrderProduct{
			OrderID:   orderID,
			ProductID: int64(product.ID),
			Quantity:  product.Quantity,
		}

		_, err = op.Insert()
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}
	s.Values["Cart"] = &models.CartCookie{Products: map[int]int{}, Locked: false}
	s.Values["Order"] = &models.OrderCookie{}

	// Saves the cookie and checks for errors
	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// TODO: send email with the invoice
	return Redirect(w, r, "/orders")
}

// checkoutPOST handles the POST request for /checkout page
func checkoutPOST(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	if !s.IsLoggedIn() {
		return http.StatusUnauthorized, nil
	}

	// Parses the form and checks for errors
	err := r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	cart, err := s.GetCart()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	order := models.OrderCookie{}
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
	fmt.Println(promocode)

	if promocode != "" {
		// Gets the promocode and checks for errors
		var generic models.Generic
		generic, err = models.GetPromocodeByCode(r.FormValue("promocode"))
		if err == sql.ErrNoRows {
			return http.StatusNotFound, nil
		}

		if err != nil {
			return http.StatusInternalServerError, err
		}

		promo := generic.(*models.Promocode)
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

	_, err = c.GetAccessToken()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	amount := paypalsdk.Amount{
		Total:    displayCents(order.Total),
		Currency: "EUR",
	}

	p, err := c.CreateDirectPaypalPayment(
		amount,
		BaseAddress+"/checkout/confirm",
		BaseAddress+"/checkout/cancel",
		"oi", // TODO: alterar esta description
	)

	if err != nil || p.ID == "" {
		return http.StatusInternalServerError, err
	}

	return Redirect(w, r, p.Links[1].Href)
}

// ValidatePromocode validates a promocode and returns the discount amount
// if it exists.
func ValidatePromocode(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	if !s.IsLoggedIn() {
		return http.StatusUnauthorized, errNotLoggedIn
	}

	code := new(bytes.Buffer)
	code.ReadFrom(r.Body)

	generic, err := models.GetPromocodeByCode(string(code.Bytes()))
	if err == sql.ErrNoRows {
		return http.StatusNotFound, nil
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}

	promocode := generic.(*models.Promocode)

	if time.Now().Unix() > promocode.Expires.Unix() {
		return http.StatusNotFound, nil
	}

	res := map[string]interface{}{}
	res["Discount"] = promocode.Discount
	res["Percentage"] = promocode.Percentage

	marsh, err := json.MarshalIndent(res, "", "")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "applicaion/json; charset=utf-8")
	if _, err := w.Write(marsh); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
