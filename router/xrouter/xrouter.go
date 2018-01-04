package xrouter

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/urfave/negroni"

	"xway/context"
	en "xway/engine"
	"xway/enum"
	xerror "xway/error"
)

// Router ...
type Router struct {
	// snp       *en.Snapshot
	frontendMap map[string]*en.Frontend
	frontends   []*en.Frontend
	// frontendMapTemp map[string]*en.Frontend
}

func (rt *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// 处理路由匹配
	fmt.Printf("[MW:xrouter] -> url router for: r.Host %v, r.URL %v\n", r.Host, r.URL)
	match, fe := rt.IsValid(r)
	if !match {
		DefaultNotFound(rw, r)
		return
	}
	// TODO: match中间件处理
	if fe == nil {
		fmt.Printf("match frontend %+v\n", fe)
	}

	next(rw, r)
}

// Remove ...
func (rt *Router) Remove(f interface{}) error {
	var frontends []*en.Frontend
	rid := f.(string)
	fmt.Printf("[xrouter.Remove] frontend %v\n", rid)
	delete(rt.frontendMap, rid)
	for _, v := range rt.frontendMap {
		fmt.Printf("%v\n", v)
		frontends = append(frontends, v)
	}
	rt.frontends = frontends
	return nil
}

// Handle ...
func (rt *Router) Handle(f interface{}) error {
	// TODO: add/update
	var frontends []*en.Frontend
	fr := f.(en.Frontend)
	fmt.Printf("[xrouter.Handle] frontend %v\n", fr)
	rt.frontendMap[fr.RouteId] = &fr
	for _, v := range rt.frontendMap {
		fmt.Printf("%v\n", v)
		frontends = append(frontends, v)
	}
	rt.frontends = frontends
	return nil
}

// frontendSlice 排序
type frontendSlice []en.Frontend

func (fs frontendSlice) Len() int {
	return len(fs)
}

func (fs frontendSlice) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (fs frontendSlice) Less(i, j int) bool {
	// 降序
	return len(strings.Split(strings.Trim(fs[j].RouteUrl, "/"), "/")) < len(strings.Split(strings.Trim(fs[i].RouteUrl, "/"), "/"))
}

// IsValid ...
func (rt *Router) IsValid(r *http.Request) (bool, interface{}) {
	if !strings.HasPrefix(r.URL.Path, "/gateway/") {
		return false, nil
	}

	xwayCtx := xwaycontext.DefaultXWayContext(r.Context())

	forwardURL := strings.Replace(r.URL.Path, "/gateway/", "/", 1)
	forwardURL = strings.ToLower(strings.TrimRight(forwardURL, "/"))
	xwayCtx.Map["forwardURL"] = forwardURL // 传递的forwardURL末尾不带"/"

	forwardURL += "/"
	var matchers []en.Frontend
	// TODO: 优化匹配逻辑
	for _, v := range rt.frontends {
		rurl := strings.ToLower(strings.TrimRight(v.RouteUrl, "/")) + "/"
		if v.DomainHost == r.Host && strings.HasPrefix(forwardURL, rurl) {
			matchers = append(matchers, *v)
		}
	}

	if len(matchers) <= 0 {
		return false, nil
	}

	// TODO: 优化, 整个路由表的排序可放在路由初始化时
	sort.Sort(frontendSlice(matchers))
	// fmt.Printf("matchers %v\n", matchers)
	res := matchers[0]
	xwayCtx.Map["matchRouteFrontend"] = &res

	return true, &res
}

// New ...
func New(snp *en.Snapshot) negroni.Handler {
	var frontends []*en.Frontend
	frontendMap := make(map[string]*en.Frontend)
	for _, v := range snp.FrontendSpecs {
		f := v.Frontend
		frontends = append(frontends, &f)
		frontendMap[v.Frontend.RouteId] = &f
	}
	return &Router{
		// snp:       snp,
		frontendMap: frontendMap,
		frontends:   frontends,
	}
}

// DefaultNotFound is an HTTP handler that returns simple 404 Not Found response.
var DefaultNotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	e := xerror.NewRequestError(enum.RetProxyError, enum.ECodeRouteNotFound, "代理路由异常")
	e.Write(w)
})
