package users

type User struct {
	FirstName        string
	LastName         string
	EmailAddress     string
	PhysicalAddress  string
	Confirmed        bool
	ConfirmationHash string
	Link             string
}
