package shopy

import "time"

// Link ...
type Link struct {
	Hash    string     `db:"hash"`
	User    int        `db:"user_id"`
	Path    string     `db:"path"`
	Used    bool       `db:"used"`
	Time    *time.Time `db:"time"`
	Expires *time.Time `db:"expires"`
}

// IsValid returns if the link is still valid and not used
func (l Link) IsValid() bool {
	return l.Expires.Unix() < time.Now().Unix() && !l.Used
}

// LinkService ...
type LinkService interface {
	Get(hash string) (*Link, error)
	Gets(first, limit int, order string) ([]*Link, error)

	Create(l *Link) error
	Update(l *Link, fields ...string) error
	Delete(hash string) error
}
