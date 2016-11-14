package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/logpacker/PayPal-Go-SDK"
	"github.com/upframe/fest"
)

// TODO:

// CheckoutCancelHandler ...
type CheckoutCancelHandler struct {
	UserService fest.UserService
}

func (h *CheckoutCancelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		err  error
	)
	defer checkErrors(w, r, code, err)

	s := r.Context().Value("session").(*fest.Session)

	if !s.IsLoggedIn() {
		if r.Method == http.MethodGet {
			code, err = Redirect(w, r, "/login")
			return
		}

		code, err = http.StatusUnauthorized, fest.ErrNotLoggedIn
		return
	}

	if r.Method != http.MethodGet {
		code = http.StatusNotImplemented
		return
	}

	cart := s.Values["Cart"].(fest.CartCookie)
	cart.Locked = false

	s.Values["Order"] = &fest.OrderCookie{}
	s.Values["Cart"] = cart

	err = s.Save(r, w)
	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	code, err = Redirect(w, r, "/cart")
}

// CheckoutConfirmHandler ...
type CheckoutConfirmHandler struct {
	UserService    fest.UserService
	ProductService fest.ProductService
	OrderService   fest.OrderService
}

func (h *CheckoutConfirmHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		err  error
	)
	defer checkErrors(w, r, code, err)

	s := r.Context().Value("session").(*fest.Session)

	if !s.IsLoggedIn() {
		if r.Method == http.MethodGet {
			code, err = Redirect(w, r, "/login")
			return
		}

		code, err = http.StatusUnauthorized, fest.ErrNotLoggedIn
		return
	}

	if r.Method != http.MethodGet {
		code = http.StatusNotImplemented
		return
	}

	paymentID := r.URL.Query().Get("paymentId")
	payerID := r.URL.Query().Get("PayerID")

	if paymentID == "" || payerID == "" {
		code, err = http.StatusBadRequest, nil
		return
	}

	_, err = fest.PayPal.GetAccessToken()
	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	executeResult, err := fest.PayPal.ExecuteApprovedPayment(paymentID, payerID)
	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	cart, err := s.GetCart(h.ProductService)
	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	order := s.Values["Order"].(fest.OrderCookie)

	o := &fest.Order{
		UserID:      s.User.ID,
		PayPalID:    paymentID,
		Value:       order.Total,
		Status:      executeResult.State,
		PromocodeID: fest.NullInt64{},
	}

	if order.Promocode.Code == "" {
		o.PromocodeID.Valid = false
	} else {
		o.PromocodeID.Valid = true
		o.PromocodeID.Int64 = int64(order.Promocode.ID)
	}

	err = h.OrderService.Create(o)
	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	o.Products = []*fest.OrderProduct{}

	for _, product := range cart.Products {
		o.Products = append(o.Products, &fest.OrderProduct{
			OrderID:   o.ID,
			ProductID: product.ID,
			Quantity:  product.Quantity,
		})
	}

	err = h.OrderService.AddProducts(o)
	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	s.Values["Cart"] = &fest.CartCookie{Products: map[int]int{}, Locked: false}
	s.Values["Order"] = &fest.OrderCookie{}

	// Saves the cookie and checks for errors
	err = s.Save(r, w)
	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	// TODO: send email with the invoice
	code, err = Redirect(w, r, "/orders")
}

// CheckoutHandler ...
type CheckoutHandler struct {
	UserService      fest.UserService
	ProductService   fest.ProductService
	PromocodeService fest.PromocodeService
}

func (h *CheckoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
func (h *CheckoutHandler) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	if !s.IsLoggedIn() {
		return Redirect(w, r, "/login")
	}

	cart, err := s.GetCart(h.ProductService)
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

	return RenderHTML(w, s, cart, "checkout")
}

// POST ...
func (h *CheckoutHandler) POST(w http.ResponseWriter, r *http.Request) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	if !s.IsLoggedIn() {
		return http.StatusUnauthorized, fest.ErrNotLoggedIn
	}

	// Parses the form and checks for errors
	err := r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	cart, err := s.GetCart(h.ProductService)
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
	fmt.Println(promocode)

	if promocode != "" {
		// Gets the promocode and checks for errors
		var promo *fest.Promocode
		promo, err = h.PromocodeService.GetByCode(promocode)
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

	_, err = fest.PayPal.GetAccessToken()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	amount := paypalsdk.Amount{
		Total:    displayCents(order.Total),
		Currency: "EUR",
	}

	p, err := fest.PayPal.CreateDirectPaypalPayment(
		amount,
		fest.BaseAddress+"/checkout/confirm",
		fest.BaseAddress+"/checkout/cancel",
		"Shop at Upframe Fest",
	)

	if err != nil || p.ID == "" {
		return http.StatusInternalServerError, err
	}

	return Redirect(w, r, p.Links[1].Href)
}

// ValidatePromocodeHandler ...
type ValidatePromocodeHandler struct {
	UserService      fest.UserService
	PromocodeService fest.PromocodeService
}

func (h *ValidatePromocodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		err  error
	)
	defer checkErrors(w, r, code, err)

	s := r.Context().Value("session").(*fest.Session)

	if !s.IsLoggedIn() {
		if r.Method == http.MethodGet {
			code, err = Redirect(w, r, "/login")
			return
		}

		code, err = http.StatusUnauthorized, fest.ErrNotLoggedIn
		return
	}

	if r.Method != http.MethodGet {
		code = http.StatusNotImplemented
		return
	}

	byt := new(bytes.Buffer)
	byt.ReadFrom(r.Body)

	promocode, err := h.PromocodeService.GetByCode(string(byt.Bytes()))
	if err == sql.ErrNoRows {
		code = http.StatusNotFound
		return
	}
	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	if time.Now().Unix() > promocode.Expires.Unix() {
		code = http.StatusNotFound
		return
	}

	res := map[string]interface{}{}
	res["Discount"] = promocode.Discount
	res["Percentage"] = promocode.Percentage

	marsh, err := json.MarshalIndent(res, "", "")
	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	w.Header().Set("Content-Type", "applicaion/json; charset=utf-8")
	if _, err := w.Write(marsh); err != nil {
		code = http.StatusInternalServerError
		return
	}

	code = http.StatusOK
}
