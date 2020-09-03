package routes

import (
	"net/http"
	"net/url"
	"time"
)

func HandleLogout(w http.ResponseWriter, r *http.Request, params url.Values) {
	sessionToken, _ := r.Cookie("session_token")
	Cache.DeleteToken(sessionToken.Value)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
