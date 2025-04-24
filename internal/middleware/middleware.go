package middleware

import (
	"iotstarter/internal/auth"
	"log"
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.Handler

func LoadMiddleware() Middleware {
	return createMiddlewareStack(
		loggingMiddleware,
		authMiddleware,
	)
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
		return true
	}
	return false
}

func isAuthed(cookieVal string) bool {
	return cookieVal == "superSecret"
}

func authMiddleware(next http.Handler) http.Handler {
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
		if !isAuthed(cookieVal) {
			http.Error(w, "No permissions to access this resource", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
