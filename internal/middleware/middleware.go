package middleware

import (
	"context"
	"iotstarter/internal/auth"
	"iotstarter/internal/config"
	"iotstarter/internal/store"
	"log"
	"net/http"
	"strings"
	"time"
)

type Middleware func(http.Handler) http.Handler

func LoadMiddleware(s store.SessionStore) Middleware {
	h := newHandler(s)
	return createMiddlewareStack(
		loggingMiddleware,
		h.authMiddleware,
	)
}

type Handler struct {
	store store.SessionStore
}

func newHandler(store store.SessionStore) Handler {
	return Handler{store: store}
}

func createMiddlewareStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func isPublicPath(path string) bool {
	if path == "/" {
		return true
	}
	publicPrefixes := []string{
		"/login",
		"/register",
		"/static",
		"/favicon.ico",
		"/api/auth/login",
		"/api/auth/register",
	}
	for _, publicPrefix := range publicPrefixes {
		if strings.HasPrefix(path, publicPrefix) {
			return true
		}
	}
	return false
}

func isSavingMeasurement(r *http.Request) bool {
	if r.URL.Path == "/api/measurements" && r.Method == "POST" {
		return auth.IsAuthedToken(r.Header.Get("x-api-key"))
	}
	return false
}

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		if isSavingMeasurement(r) {
			next.ServeHTTP(w, r)
			return
		}

		cookieVal, err := auth.GetCookieValue(r)
		if err != nil {
			http.Error(w, "No session value", http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		user, err := h.store.GetUserFromToken(ctx, cookieVal)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unauthenticated", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(r.Context(), config.UserKey, user)
		r = r.WithContext(ctx)

		log.Println("User set in request context", user)

		next.ServeHTTP(w, r)
	})
}
