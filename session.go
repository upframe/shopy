package upframe

import "time"

type Session struct {
	User         *User
	SessionKey   string
	LoginTime    *time.Time
	LastSeenTime *time.Time
}
