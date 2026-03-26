package router

import (
	"strings"

	"flexgate/internal/config"
)

type Router struct {
	routes []config.Route
}

func NewRouter(routes []config.Route) *Router {
	return &Router{routes: routes}
}

func (r *Router) Match(path string) *config.Route {
	for _, route := range r.routes {
		if strings.HasPrefix(path, route.PathPrefix) {
			return &route
		}
	}
	return nil
}
