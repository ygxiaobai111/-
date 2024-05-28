package net

import (
	"log"
	"strings"
	"sync"
)

//account login||logout

type HandlerFunc func(req *WsMsgReq, rsp *WsMsgRsp)

type MiddlewareFunc func(handlerFunc HandlerFunc) HandlerFunc

type group struct {
	mutex         sync.RWMutex
	prefix        string
	handlerMap    map[string]HandlerFunc
	middlewareMap map[string][]MiddlewareFunc //路由中间件
	middlewares   []MiddlewareFunc            //全局中间件
}

// 路由组中添加方法
func (g *group) AddRouter(name string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.handlerMap[name] = handlerFunc

	g.middlewareMap[name] = middlewares

}

// Use 添加全局中间件
func (g *group) Use(middlewares ...MiddlewareFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// 添加路由组
func (r *Router) Group(prefix string) *group {
	g := &group{
		prefix:        prefix,
		handlerMap:    make(map[string]HandlerFunc),
		middlewareMap: make(map[string][]MiddlewareFunc),
		middlewares:   make([]MiddlewareFunc, 0),
	}
	r.group = append(r.group, g)

	return g
}
func (g *group) exec(name string, req *WsMsgReq, rsp *WsMsgRsp) {

	h, ok := g.handlerMap[name]
	if !ok {
		//推断是否为网关
		h, ok = g.handlerMap["*"]
		if !ok {
			log.Println("路由未定义")
		}
	}
	if ok {
		//中间件 执行路由之前 需要执行中间件代码
		for i := 0; i < len(g.middlewares); i++ {
			h = g.middlewares[i](h)
		}
		mm, ok := g.middlewareMap[name]
		if ok {
			for i := 0; i < len(mm); i++ {
				h = mm[i](h)
			}
		}
		h(req, rsp)
	}
}

type Router struct {
	group []*group
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Run(req *WsMsgReq, rsp *WsMsgRsp) {
	//req.Body.Name 路径 登录业务 account.login （account组标识）login 路由标识
	strs := strings.Split(req.Body.Name, ".")
	prefix := ""
	name := ""
	if len(strs) == 2 {
		prefix = strs[0]
		name = strs[1]
	}

	for _, g := range r.group {
		log.Println(g.prefix, " ", prefix)
		if g.prefix == prefix {
			g.exec(name, req, rsp)
		} else if g.prefix == "*" {
			g.exec(name, req, rsp)
		}
	}
}
