package models

import "time"

const (
	deleteExpiration  = 60 * 60 * 2       // 2 Hours in seconds
	confirmExpiration = 60 * 60 * 24 * 20 // 20 Days in seconds
	resetExpiration   = 60 * 60 * 2       // 2 Hours in seconds
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
