package http

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/upframe/fest"
)

// TODO?

const itemsPerPage = 10

type adminTable struct {
	Items        interface{}
	HasPrevious  bool
	LinkPrevious string
	HasNext      bool
	LinkNext     string
}

// AdminGet ...
func AdminGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)
	return Render(w, c, s, nil, "admin/home")
}

// AdminRedirect ...
func AdminRedirect(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	return Redirect(w, r, "/admin/"+mux.Vars(r)["category"]+"/1")
}

// AdminNew ...
func AdminNew(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)
	category := mux.Vars(r)["category"]
	return Render(w, c, s, nil, "admin/"+category+"-new", "admin/"+category+"-form")
}

// AdminListing ...
func AdminListing(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	category := mux.Vars(r)["category"]

	// Gets the number of the page.
	page, err := strconv.Atoi(mux.Vars(r)["page"])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	data := &adminTable{}
	// Calculates the offset and gets the item.
	offset := (page - 1) * itemsPerPage

	switch category {
	case "orders":
		data.Items, err = c.Services.Order.Gets(offset, itemsPerPage, "ID")
	case "products":
		data.Items, err = c.Services.Product.Gets(offset, itemsPerPage, "ID")
	case "promocodes":
		data.Items, err = c.Services.Promocode.Gets(offset, itemsPerPage, "ID")
	case "users":
		data.Items, err = c.Services.User.Gets(offset, itemsPerPage, "ID")
	default:
		// This mustn't happen!
		return http.StatusInternalServerError, nil
	}

	// Checks if there are any item. If we're in the first page, show
	// it anyway so we're able to create new item.
	if err == sql.ErrNoRows && page != 1 {
		return http.StatusNotFound, err
	}

	// Checks for other errors.
	if err != nil {
		return http.StatusInternalServerError, err
	}

	data.HasPrevious = page != 1
	data.LinkPrevious = "/admin/" + category + "/" + (strconv.Itoa(page - 1))

	var total int
	switch category {
	case "orders":
		total, err = c.Services.Order.Total()
	case "products":
		total, err = c.Services.Product.Total()
	case "promocodes":
		total, err = c.Services.Promocode.Total()
	case "users":
		total, err = c.Services.User.Total()
	default:
		// This mustn't happen!
		return http.StatusInternalServerError, nil
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	data.HasNext = total > (itemsPerPage * page)
	data.LinkNext = "/admin/" + category + "/" + strconv.Itoa(page+1)

	// Show the page with the table.
	return Render(w, c, s, data, "admin/"+category, "admin/"+category+"-form")
}
