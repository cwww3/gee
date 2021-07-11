package gee

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"runtime/debug"
	"strings"
)

type Engine struct {
	router *router
	*RouterGroup
	groups        []*RouterGroup
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	engine.Use(Recover())
	return engine
}

type HandlerFunc func(c *Context)

func Recover() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("err=", err)
				fmt.Printf("%s\n", debug.Stack())
				c.String(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		c.Next()
	}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, r)
	c.handlers = middlewares
	c.engine = e
	e.router.handle(c)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.router.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.router.addRoute("POST", pattern, handler)
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

func (group *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	childGroup := &RouterGroup{
		prefix:      group.prefix + prefix,
		middlewares: nil,
		parent:      group,
		engine:      engine,
	}
	engine.groups = append(group.engine.groups, childGroup)
	return childGroup
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}
