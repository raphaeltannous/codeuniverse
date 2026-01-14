package handlersutils

import "net/http"

func WriteJwtCookie(w http.ResponseWriter, jwt string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    jwt,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
