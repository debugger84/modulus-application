package application

import "net/http"

type Routes struct {
	routes map[string]RouteInfo
}

func NewRoutes() *Routes {
	return &Routes{routes: make(map[string]RouteInfo)}
}

func (r *Routes) Get(name string, path string, handler http.HandlerFunc) {
	r.routes[name] = RouteInfo{
		handler: handler,
		method:  http.MethodGet,
		path:    path,
	}
}

func (r *Routes) Post(name string, path string, handler http.HandlerFunc) {
	r.routes[name] = RouteInfo{
		handler: handler,
		method:  http.MethodPost,
		path:    path,
	}
}

func (r *Routes) Delete(name string, path string, handler http.HandlerFunc) {
	r.routes[name] = RouteInfo{
		handler: handler,
		method:  http.MethodDelete,
		path:    path,
	}
}

func (r *Routes) Put(name string, path string, handler http.HandlerFunc) {
	r.routes[name] = RouteInfo{
		handler: handler,
		method:  http.MethodPut,
		path:    path,
	}
}

func (r *Routes) Options(name string, path string, handler http.HandlerFunc) {
	r.routes[name] = RouteInfo{
		handler: handler,
		method:  http.MethodOptions,
		path:    path,
	}
}

func (r *Routes) AddFromRoutes(routes *Routes) {
	for name, info := range routes.routes {
		r.routes[name] = info
	}
}

func (r *Routes) GetRoutesInfo() []RouteInfo {
	result := make([]RouteInfo, 0, len(r.routes))
	for _, info := range r.routes {
		result = append(result, info)
	}

	return result
}
