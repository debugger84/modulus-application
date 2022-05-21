package application

import "net/http"

type Router interface {
	AddRoutes(routes []RouteInfo)
	Run() error
}

type RouteInfo struct {
	method  string
	path    string
	handler http.HandlerFunc
}

func NewRouteInfo(method string, path string, handler http.HandlerFunc) *RouteInfo {
	return &RouteInfo{method: method, path: path, handler: handler}
}

func (r RouteInfo) Handler() http.HandlerFunc {
	return r.handler
}

func (r RouteInfo) Method() string {
	return r.method
}

func (r RouteInfo) Path() string {
	return r.path
}
