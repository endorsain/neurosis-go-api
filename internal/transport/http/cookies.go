package httptransport

import (
	"net/http"
	"strings"
	"time"

	"github.com/endorsain/neurosis-go-api/internal/config"
)

func SetRefreshTokenCookie(w http.ResponseWriter, value string, expiresAt time.Time, cfg config.RefreshCookieConfig) {
	sameSite := http.SameSiteLaxMode
	switch strings.ToLower(cfg.SameSite) {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cfg.Name,
		Value:    value,
		Path:     cfg.Path,
		Domain:   cfg.Domain,
		Secure:   cfg.Secure,
		HttpOnly: cfg.HTTPOnly,
		SameSite: sameSite,
		Expires:  expiresAt,
	})
}
