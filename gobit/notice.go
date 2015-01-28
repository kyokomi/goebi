package gobit

import (
	"fmt"
	"net/http"
	"runtime"
	"os"
)

type Notice struct {
	Notifier notifier               `json:"notifier"`
	Errors   []Error                `json:"errors"`
	Env      map[string]interface{} `json:"environment"`
	Session  map[string]interface{} `json:"session"`
	Params   map[string]interface{} `json:"params"`
}

type Error struct {
	Type      string       `json:"type"`
	Message   string       `json:"message"`
	Backtrace []StackFrame `json:"backtrace"`
}

type notifier struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	URL     string `json:"url"`
}

type StackFrame struct {
	File string `json:"file"`
	Line int    `json:"line"`
	Func string `json:"function"`
}

// TODO: 自動生成したい
var app = notifier{
	Version: "0.1.0",
	Name: "gobit",
	URL: "https://github.com/kyokomi/gobit",
}

func NewNoticeWithFilter(err interface{}, filter stackFilterFunc, req *http.Request) *Notice {
	return newNotice(err, stackTrace(filter), req)
}

func NewNotice(err interface{}, req *http.Request) *Notice {
	return newNotice(err, defaultStackTrace(), req)
}

// TODO: 未整理
func newNotice(err interface{}, stack []StackFrame, req *http.Request) *Notice {
	
	notice := &Notice{
		Notifier: app,
		Errors: []Error{
			Error{
				Type:      fmt.Sprintf("%T", err),
				Message:   fmt.Sprint(err),
				Backtrace: stack,
			},
		},
		Env:     createContext(),
		Session: map[string]interface{}{},
		Params:  map[string]interface{}{},
	}

	if req == nil {
		return notice
	}
	
	notice.Env["url"] = req.URL.String()
	
	if ua := req.Header.Get("User-Agent"); ua != "" {
		notice.Env["userAgent"] = ua
	}

	for k, v := range req.Header {
		if len(v) == 1 {
			notice.Env[k] = v[0]
		} else {
			notice.Env[k] = v
		}
	}

	// TODO: jsonのParamsがとれない　いずれ...
	if err := req.ParseForm(); err == nil {
		for k, v := range req.Form {
			if len(v) == 1 {
				notice.Params[k] = v[0]
			} else {
				notice.Params[k] = v
			}
		}
	}

	return notice
}

// TODO: あとでstructにするかも
func createContext() map[string]interface{} {
	context := map[string]interface{}{
		"language":     runtime.Version(),
		"os":           runtime.GOOS,
		"architecture": runtime.GOARCH,
	}
	if hostname, err := os.Hostname(); err == nil {
		context["hostname"] = hostname
	}
	if wd, err := os.Getwd(); err == nil {
		context["rootDirectory"] = wd
	}

	return context
}
