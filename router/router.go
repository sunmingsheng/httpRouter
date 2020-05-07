package router

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"net/http"
)

var nextFlagId string

func init() {
	nextFlagId, _ = uuid.GenerateUUID()
}

//http服务调度器
type Mux struct {
	serveMux *http.ServeMux
	Router *Group
}

//路由分组
type Group struct {
	serveMux *http.ServeMux
	groupPath string
	handlers []middlewareFunc
}

//定义中间件函数类型
type middlewareFunc func (ctx *context.Context) http.Handler

//生成一个新的调度器
func NewMux() *Mux {
	serveMux := http.NewServeMux()
	mux := &Mux{
		serveMux: serveMux,
	}
	//添加默认分组
	mux.Router = mux.AddDefaultGroup("/", defaultMiddleware)
	return mux
}

//启动http服务
func (Mux *Mux) StartMux(port string) error {
	if err := http.ListenAndServe(port, Mux.serveMux); err != nil {
		return err
	}
	return nil
}

//添加默认路由分组
func (Mux *Mux) AddDefaultGroup(groupPath string, handlers ...middlewareFunc) *Group {
	return &Group {
		serveMux: Mux.serveMux,
		groupPath: groupPath,
		handlers:  handlers,
	}
}

//添加路由分组
func (group *Group) AddGroup(groupPath string, handlers ...middlewareFunc) *Group {
	return &Group{
		serveMux: group.serveMux,
		groupPath: group.groupPath + groupPath,
		handlers:  append(group.handlers, handlers...),
	}
}

//为分组添加具体路由
func (group *Group) Any(path string,f func (w http.ResponseWriter, r *http.Request)) {
	h := requestHandler(group.handlers, f, "Any")
	group.serveMux.HandleFunc(group.groupPath + path, h)
}

//为分组添加具体路由
func (group *Group) Get(path string,f func (w http.ResponseWriter, r *http.Request)) {
	h := requestHandler(group.handlers, f, "Get")
	group.serveMux.HandleFunc(group.groupPath + path, h)
}

//为分组添加具体路由
func (group *Group) Post(path string,f func (w http.ResponseWriter, r *http.Request)) {
	h := requestHandler(group.handlers, f, "Post")
	group.serveMux.HandleFunc(group.groupPath + path, h)
}

//处理函数继续往下执行
func Next(ctx *context.Context) {
	*ctx = context.WithValue(*ctx, nextFlagId, true)
}

//处理请求
func requestHandler(midFuncSlices []middlewareFunc, f func (w http.ResponseWriter, r *http.Request), acceptMethod string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch acceptMethod {
		case "Any":
		case "Post":
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
		case "Get":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
		}

		ctx := context.Background()
		for _, middlewareFunc := range midFuncSlices {
			ctx = context.WithValue(ctx, nextFlagId, false)
			h := middlewareFunc(&ctx)
			h.ServeHTTP(w, r)
			if ctx.Value(nextFlagId).(bool) == false {
				return
			}
			r = r.WithContext(ctx)
		}
		f(w, r)
	}
}

//定义基础中间件
func defaultMiddleware(ctx *context.Context) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		Next(ctx)
	})
}




