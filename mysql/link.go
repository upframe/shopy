package mysql

import (
	"github.com/jmoiron/sqlx"
	"github.com/bruhs/shopy"
)

// LinkService ...
type LinkService struct {
	DB *sqlx.DB
}

var linkMap = map[string]string{
	"Hash":    "hash",
	"User":    "user_id",
	"Path":    "path",
	"Used":    "used",
	"Time":    "time",
	"Expires": "expires",
}

// Get ...
func (s *LinkService) Get(hash string) (*shopy.Link, error) {
	link := &shopy.Link{}
	err := s.DB.Get(link, "SELECT * FROM links WHERE hash=?", hash)

	return link, err
}

// Gets ...
func (s *LinkService) Gets(first, limit int, order string) ([]*shopy.Link, error) {
	links := []*shopy.Link{}
	var err error

	order = fieldsToColumns(linkMap, order)[0]

	if limit == 0 {
		err = s.DB.Select(&links, "SELECT * FROM links ORDER BY ?", order)
	} else {
		err = s.DB.Select(&links, "SELECT * FROM links ORDER BY ? LIMIT ? OFFSET ?", order, limit, first)
	}

	return links, err
}

// Create ...
func (s *LinkService) Create(l *shopy.Link) error {
	_, err := s.DB.NamedExec(insertQuery("links", getAllColumns(linkMap)), l)
	return err
}

// Update ...
func (s *LinkService) Update(l *shopy.Link, fields ...string) error {
	_, err := s.DB.NamedExec(updateQuery("links", "hash", fieldsToColumns(linkMap, fields...)), l)
	return err
}

// Delete ...
func (s *LinkService) Delete(hash string) error {
	l, err := s.Get(hash)
	if err != nil {
		return err
	}

	l.Used = true
	return s.Update(l, "Used")
}
