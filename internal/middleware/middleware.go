package middleware

import (
	"context"
	"iotstarter/internal/auth"
	"iotstarter/internal/model"
	"iotstarter/internal/security"
	"iotstarter/internal/session"
	"log"
	"net/http"
	"strings"
	"time"
)

type Middleware func(http.Handler) http.Handler

func LoadMiddleware(s *session.Service) Middleware {
	h := newHandler(s)
	return createMiddlewareStack(
		loggingMiddleware,
		h.authMiddleware,
	)
}

type Handler struct {
	s *session.Service
}

func newHandler(s *session.Service) Handler {
	return Handler{s: s}
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
		"/register",
		"/login",
		"/logout",
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
		return security.IsAuthedToken(r.Header.Get("x-api-key"))
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

		cookieVal, err := session.GetCookieValue(r)
		if err != nil {
			http.Error(w, "No session value", http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		token := model.SessionToken(cookieVal)

		user, err := h.s.GetUserFromToken(ctx, token)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unauthenticated", http.StatusUnauthorized)
			return
		}

		log.Println("User in session", user)

		ctx = auth.WithUser(r.Context(), user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
