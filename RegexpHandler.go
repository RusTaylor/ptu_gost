package main

import (
	"net/http"
	"regexp"
	"time"
)

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexHandler struct {
	routes []*route
}

func (h *RegexHandler) HttpHandler(pattern *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

func (h *RegexHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *RegexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, _ := r.Cookie("token")
	if token == nil {
		cookie := &http.Cookie{
			Name:    "token",
			Value:   "token_test",
			Expires: time.Now().Add(60 * time.Minute),
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}
