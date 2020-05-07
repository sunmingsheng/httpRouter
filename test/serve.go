package main

import (
	"context"
	"fmt"
	"httpRouter"
	"net/http"
)

func main() {
	mux := httpRouter.NewMux()
	r := mux.Router
	r.Get("test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})
	g := r.AddGroup("admin", middleware2, middleware1)
	g.Any("/test", func (w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Context().Value("age"))
		fmt.Println(r.Context().Value("sex"))
		w.Write([]byte("inner..."))
	})
	g2 := g.AddGroup("/2", middleware2)
	g2.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("2/test"))
	})

	if err := mux.StartMux(":8080"); err != nil {
		panic(err)
	}
}

//定义中间件
func middleware1(ctx *context.Context) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		if r.Form.Get("age") == "" {
			*ctx = context.WithValue(*ctx, "age", 22)
		} else {
			*ctx = context.WithValue(*ctx, "age", r.Form.Get("age"))
		}
		w.Write([]byte("mid1..."))
		httpRouter.Next(ctx)
	})
}

func middleware2(ctx *context.Context) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		if r.Form.Get("sex") == "" {
			*ctx = context.WithValue(*ctx, "sex","男")
		} else {
			*ctx = context.WithValue(*ctx, "sex",r.Form.Get("sex"))
		}
		w.Write([]byte("mid2..."))
		httpRouter.Next(ctx)
	})
}

