package pages

import (
	"net/http"

	"github.com/upframe/fest/models"
)

// AdminUsersGET is
func AdminUsersGET(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	return AdminGenericGET(w, r, s, "users", models.GetUsers)
}

// AdminUsersPOST is
func AdminUsersPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	return AdminGenericPOST(w, r, new(models.User))
}

// AdminUsersDELETE is
func AdminUsersDELETE(w http.ResponseWriter, r *http.Request) (int, error) {
	return AdminGenericDELETE(w, r, "users", models.GetUserByID)
}

// AdminUsersPUT is
func AdminUsersPUT(w http.ResponseWriter, r *http.Request) (int, error) {
	// In Users, we do not update: password_hash, password_salt, referral,
	// nor referrer
	return AdminGenericPUT(w, r, new(models.User), "id",
		"first_name",
		"last_name",
		"email",
		"address",
		"invites",
		"credit",
		"confirmed",
		"admin",
		"deactivated")
}
