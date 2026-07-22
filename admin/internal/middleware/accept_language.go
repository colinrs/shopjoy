package middleware

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/rest"
)

// AcceptLanguageHeader is the canonical HTTP header name carrying locale
// preference. We read it verbatim (including q-values and trailing tags) so
// downstream resolvers can pick the right parsing strategy.
const AcceptLanguageHeader = "Accept-Language"

// AcceptLanguage injects the raw Accept-Language header into ctx under
// contextx.acceptLanguageKey. Logic / application layers can then call
// contextx.GetAcceptLanguage(ctx) to recover the original header string.
//
// Empty / missing header → ctx value is "" (no locale signal); resolvers
// must treat that as "use default fallback" rather than an error.
//
// Usage: register on any route that resolves localized strings (e.g.
// shipping zones). Cheap enough to leave on every authenticated route if
// needed, but no harm in keeping it scoped to the routes that consume it.
func AcceptLanguage() rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := contextx.SetAcceptLanguage(r.Context(), r.Header.Get(AcceptLanguageHeader))
			next(w, r.WithContext(ctx))
		}
	}
}