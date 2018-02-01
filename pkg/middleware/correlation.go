package middleware

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/frankgreco/demogods/pkg/log"
)

// Correlation is a middleware that injects a correlation id into
// the request's context. This middleware is most effective if
// execeted before other middleware.
func Correlation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := log.NewContext(r.Context(), zap.Stringer("correlation_id", uuid.NewV4()))
		log.WithContext(ctx).Info("established new correlation id for this request")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
