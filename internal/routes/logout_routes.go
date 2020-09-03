package routes

import (
	"net/http"
	"net/url"
	"time"

	uuid "github.com/satori/go.uuid"
)

func HandleLogout(w http.ResponseWriter, r *http.Request, params url.Values) {
	sessionToken, _ := uuid.NewV4()
	Cache.DeleteToken(sessionToken.String())

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
