package models

import "time"

// Link is an object that holds an hash, the corresponding user, an action,
// the moment it was created and an expiration date.
//
// These links are used to send confirmation emails, delete confirmation emails
// and reset emails.
type Link struct {
	Hash    string     `db:"hash"`
	User    int        `db:"user_id"`
	Path    string     `db:"path"`
	Used    bool       `db:"used"`
	Time    *time.Time `db:"time"`
	Expires *time.Time `db:"expires"`
}

// Update updates the current link object in the database
func (l Link) Update(fields ...string) error {
	_, err := db.NamedExec(updateQuery("links", "hash", fields), l)
	return err
}

// Insert inserts the current Link struct into the database and returns an error
// if something goes wrong.
func (l Link) Insert() error {
	_, err := db.NamedExec(
		`INSERT INTO links
		            (hash,
		             user_id,
		             path,
		             used,
		             time,
		             expires)
		VALUES      (:hash,
                     :user_id,
                     :path,
                     :used,
                     :time,
                     :expires)`, l)

	return err
}

// IsValid returns if the link is still valid and not used
func (l Link) IsValid() bool {
	return l.Expires.Unix() < time.Now().Unix() && !l.Used
}

// GetLinkByHash returns a link using its hash and an error if your database sucks or your code sucks
func GetLinkByHash(hash string) (*Link, error) {
	link := &Link{}
	err := db.Get(link, "SELECT * FROM links WHERE hash=?", hash)

	return link, err
}
