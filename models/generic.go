package models

// Generic that is used to simplify code instead of having functions for Promocodes, Orders and Products
type Generic interface {
	Insert() error
	Update(fields ...string) error
	Deactivate() error
}

type GetGeneric func(int) (Generic, error)
type GetGenerics func(int, int) ([]Generic, error)
