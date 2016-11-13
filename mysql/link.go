package mysql

import (
	"github.com/upframe/fest"
)

// LinkService ...
type LinkService struct{}

var linkMap = map[string]string{
	"Hash":    "hash",
	"User":    "user_id",
	"Path":    "path",
	"Used":    "used",
	"Time":    "time",
	"Expires": "expires",
}

// Get ...
func (s *LinkService) Get(hash string) (*fest.Link, error) {
	link := &fest.Link{}
	err := db.Get(link, "SELECT * FROM links WHERE hash=?", hash)

	return link, err
}

// Gets ...
func (s *LinkService) Gets(first, limit int, order string) ([]*fest.Link, error) {
	links := []*fest.Link{}
	var err error

	if limit == 0 {
		err = db.Select(&links, "SELECT * FROM links ORDER BY ?", order)
	} else {
		err = db.Select(&links, "SELECT * FROM links ORDER BY ? LIMIT ? OFFSET ?", order, limit, first)
	}

	return links, err
}

// Create ...
func (s *LinkService) Create(l *fest.Link) error {
	_, err := db.NamedExec(insertQuery("links", getAllColumns(linkMap)), l)
	return err
}

// Update ...
func (s *LinkService) Update(l *fest.Link, fields ...string) error {
	_, err := db.NamedExec(updateQuery("links", "hash", fieldsToColumns(linkMap, fields...)), l)
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
