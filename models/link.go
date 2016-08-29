package models

import "time"

const (
	DeleteExpiration  = 60 * 60 * 2       // 2 Hours in seconds
	ConfirmExpiration = 60 * 60 * 24 * 20 // 20 Days in seconds
	ResetExpiration   = 60 * 60 * 2       // 2 Hours in seconds
)

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
	_, err := db.NamedExec(updateQuery("links", "hash", l.Hash, fields), l)
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

// NewConfirmationLink generates a new confirmation link to be used within
// confirmation emails
func NewConfirmationLink(u *User) *Link {
	now := time.Now()
	expires := now.Add(confirmExpiration)

	link := &Link{
		Path:    "/register",
		Hash:    UniqueHash(u.Email),
		User:    u.ID,
		Used:    false,
		Time:    &now,
		Expires: &expires,
	}

	return link
}
