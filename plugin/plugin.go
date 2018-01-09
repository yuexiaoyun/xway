package plugin

import (
	"fmt"
	"xway/router"

	"github.com/garyburd/redigo/redis"
	"github.com/urfave/negroni"
)

// Registry contains common obj.
type Registry struct {
	router router.Router
	// Middlewares map[string]Middleware
	middlewareSpecs map[string]*MiddlewareSpec
	redisPool       *redis.Pool
}

// MiddlewareSpec ...
// 中间件定义
type MiddlewareSpec struct {
	Type string
	MW   Middleware
}

// Middleware ...
// 中间件创建模板
type Middleware func(interface{}) negroni.Handler

// New ...
func New() *Registry {
	return &Registry{
		middlewareSpecs: make(map[string]*MiddlewareSpec),
	}
}

func (r *Registry) AddMW(mw *MiddlewareSpec) error {
	if mw == nil {
		return fmt.Errorf("middleware spec can not be nil")
	}
	if r.GetMW(mw.Type) != nil {
		return fmt.Errorf("moddleware of type %s already registered", mw.Type)
	}
	r.middlewareSpecs[mw.Type] = mw
	return nil
}

func (r *Registry) GetMW(mtype string) *MiddlewareSpec {
	return r.middlewareSpecs[mtype]
}

func (r *Registry) SetRouter(router router.Router) error {
	r.router = router
	return nil
}

func (r *Registry) GetRouter() router.Router {
	return r.router
}

func (r *Registry) SetRedisPool(p *redis.Pool) {
	r.redisPool = p
}

func (r *Registry) GetRedisPool() *redis.Pool {
	return r.redisPool
}

func (r *Registry) Close() {
	if r.redisPool != nil {
		r.redisPool.Close()
	}
}