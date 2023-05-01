package main

//go:generate vite build

import (
	"embed"
	"mime"
	"strings"

	"github.com/wasmcloud/actor-tinygo"
	"github.com/wasmcloud/interfaces/httpserver/tinygo"
	logging "github.com/wasmcloud/interfaces/logging/tinygo"
)

//go:embed dist
var f embed.FS

type UI struct {
	Logger *logging.LoggingSender
}

func main() {
	ui := UI{
		Logger: logging.NewProviderLogging(),
	}
	actor.RegisterHandlers(httpserver.HttpServerHandler(&ui))
}

func (u *UI) HandleRequest(ctx *actor.Context, req httpserver.HttpRequest) (*httpserver.HttpResponse, error) {
	ret := httpserver.HttpResponse{
		Header: make(httpserver.HeaderMap),
	}

	_ = u.Logger.WriteLog(ctx, logging.LogEntry{Level: "debug", Text: "ENDPOINT: " + req.Path})

	if req.Path == "/healthz" {
		ret.StatusCode = 200
		ret.Body = []byte("healthy!")
		return &ret, nil
	}

	path := req.Path
	if req.Path == "/" {
		path = "/index.html"
	}

	page, err := f.ReadFile("dist" + path)
	if err != nil {
		_ = u.Logger.WriteLog(ctx, logging.LogEntry{Level: "error", Text: err.Error()})
		ret.StatusCode = 404
		ret.Body = []byte("page not found")
		return &ret, nil
	}

	splitPath := strings.Split(path, ".")
	if len(splitPath) > 1 {
		ext := splitPath[len(splitPath)-1]
		ret.Header["Content-Type"] = httpserver.HeaderValues{mime.TypeByExtension("." + ext)}
		_ = u.Logger.WriteLog(ctx, logging.LogEntry{Level: "debug", Text: "DETECTED MIME TYPE: " + ret.Header["Content-Type"][0]})
	}

	ret.StatusCode = 200
	ret.Body = page

	return &ret, nil
}
