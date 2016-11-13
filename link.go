package fest

import "time"

// Link ...
type Link struct {
	Hash    string     `db:"hash"`
	User    int      `db:"user_id"`
	Path    string     `db:"path"`
	Used    bool       `db:"used"`
	Time    *time.Time `db:"time"`
	Expires *time.Time `db:"expires"`
}

// LinkService ...
type LinkService interface {
	Get(id int) (*Link, error)
	GetByHash(hash string) (*Link, error)
	Gets(first, limit int, order string) ([]*Link, error)

	Create(l *Link) error
	Update(l *Link, fields ...string) error
	Delete(id int) error
}
