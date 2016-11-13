package models

// Generic that is used to simplify code instead of having functions for Promocodes,
// Orders and Products
type Generic interface {
	Insert() (int, error)
	Update(fields ...string) error
	Deactivate() error
}

// GetGeneric is the type for GetX functions
type GetGeneric func(int) (Generic, error)

// GetGenerics is the type for GetXs functions
type GetGenerics func(int, int, string) ([]Generic, error)
