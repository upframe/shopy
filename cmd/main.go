package main

import (
	"encoding/gob"

	"github.com/upframe/fest"
)

func init() {
	// Regist types so they can be used on Cookies
	gob.Register(fest.CartCookie{})
	gob.Register(fest.OrderCookie{})
}

func main() {

}
